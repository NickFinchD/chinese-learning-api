// Command seed (re)generates lessons 2.. for the "HSK 1" course, extending
// the single hand-authored lesson 1 ("Приветствие") with programmatically
// assembled practice content: new words are introduced in topical batches
// (e.g. "Местоимения", "Числа 2") at most 3 at a time, mixed with review
// words drawn from everything introduced so far, plus vocabulary quizzes and
// grammar quizzes (from a hand-written HSK1 grammar bank). The generated
// range is followed by the hand-authored grammar-explanation lesson and its
// sentence-builder review lesson (content seeded by migrations
// 000026/000029, placed here at whatever lesson_number comes next). It is
// meant to be run once against a database that already has migrations
// 000001..000036 applied:
//
//	go run ./cmd/seed
//
// It is destructive-but-idempotent with respect to lesson numbers 2.. for
// the HSK1 course: it deletes any lessons already in that range (and their
// steps/progress, via ON DELETE CASCADE) before regenerating them, so it can
// be re-run freely to pick up changes to the topic list or grammar bank.
// Vocabulary/grammar quizzes are looked up by question text and reused
// instead of being duplicated.
//
// Grammar quizzes are gated: a grammarBank entry is only eligible once every
// word its example sentence uses has actually been introduced in an earlier
// (or the current) lesson — otherwise a learner would be quizzed on 的/吗/在
// sentences built from vocabulary they haven't seen yet. The first time a
// given grammar *construct* (e.g. "的 — possession", shared by more than one
// bank entry) is actually used, a short grammar-note step explaining it is
// inserted right before that quiz step.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
	"github.com/jackc/pgx/v5"
)

const (
	firstGeneratedLesson    = 2
	totalLessons            = 60
	wordStepsPerLesson      = 12
	bonusReviewQuizzes      = 3 // extra quizzes drawn from review words, on top of the guaranteed one-per-new-word
	grammarQuizzesPerLesson = 2
)

type wordInfo struct {
	ID          int64
	Hanzi       string
	Translation string
}

type topicGroup struct {
	name  string
	hanzi []string
}

// Topics cover every HSK1 word not already introduced by the hand-authored
// lesson 1 ("Приветствие"). Ordered so grammatically load-bearing vocabulary
// (pronouns, particles, numbers, measure words, question words) comes
// first — almost every grammar-quiz sentence leans on one of these, so
// teaching them early lets grammar quizzes unlock much sooner instead of
// piling up at the very end of the course.
var topics = []topicGroup{
	{"Местоимения", []string{"我", "你", "他", "她", "那", "这"}},
	{"Частицы и грамматика", []string{"是", "不", "有", "的", "都", "和", "了", "吗", "没", "呢", "也"}},
	{"Числа", []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}},
	{"Счётные слова", []string{"本", "分钟", "个", "号", "块", "岁", "些"}},
	{"Вопросительные слова", []string{"多少", "几", "哪", "哪儿", "谁", "什么", "怎么", "怎么样"}},
	{"Семья", []string{"爸爸", "妈妈", "儿子", "女儿", "我们", "你们"}},
	{"Глаголы и действия", []string{"吃", "喝", "看", "听", "说", "读", "写", "买", "去", "来", "做", "喜欢"}},
	{"Модальные глаголы", []string{"会", "能", "请", "想"}},
	{"Прилагательные", []string{"大", "多", "高兴", "好", "很", "冷", "漂亮", "热", "少", "太", "一点儿", "小"}},
	{"Время и даты", []string{"今天", "明天", "昨天", "现在", "年", "月", "星期", "点"}},
	{"Еда и напитки", []string{"茶", "水", "米饭", "菜", "水果", "苹果"}},
	{"Люди и профессии", []string{"工作", "老师", "朋友", "人", "同学", "先生", "小姐", "学生", "医生"}},
	{"Дом и вещи", []string{"杯子", "电脑", "电视", "电影", "东西", "家", "书", "衣服", "椅子", "桌子"}},
	{"Города и места", []string{"北京", "中国", "饭店", "商店", "学校", "医院"}},
	{"Транспорт и покупки", []string{"出租车", "飞机", "钱"}},
	{"Место и направление", []string{"后面", "里", "前面", "上", "下", "在"}},
	{"Погода и время", []string{"上午", "时候", "天气", "下午", "下雨", "中午"}},
	{"Глаголы движения", []string{"回", "开", "睡觉", "住", "坐"}},
	{"Глаголы общения", []string{"爱", "看见", "认识", "打电话", "叫"}},
	{"Вежливые фразы", []string{"不客气", "对不起", "没关系", "喂"}},
	{"Учёба", []string{"汉语", "名字", "学习", "字"}},
	{"Животные", []string{"狗", "马", "猫"}},
}

type gramOption struct {
	text    string
	correct bool
}

type gramQuestion struct {
	question string
	options  []gramOption

	// requiredHanzi lists every word the example sentence relies on — the
	// question is only eligible for a lesson once all of these have been
	// introduced (in this lesson or an earlier one).
	requiredHanzi []string

	// constructTitle names the grammar_notes row (seeded by migrations
	// 000026/000036) explaining the pattern this question tests. Several
	// entries can share a construct; its note is only inserted once, the
	// first time any entry using it is actually selected for a lesson.
	constructTitle string
}

var grammarBank = []gramQuestion{
	{"Выберите правильный порядок слов: «Я пью чай»", []gramOption{
		{"我喝茶", true}, {"我茶喝", false}, {"喝我茶", false}, {"茶我喝", false},
	}, []string{"我", "喝", "茶"}, "Порядок слов: подлежащее + время + сказуемое + дополнение"},

	{"Как правильно спросить «Это твоя книга?» (вопрос да/нет)", []gramOption{
		{"这是你的书吗?", true}, {"这是你的书呢?", false}, {"这是你的书了?", false}, {"这是你的书的?", false},
	}, []string{"这", "是", "你", "的", "书", "吗"}, "Вопрос с частицей 吗"},

	{"Выберите правильный вариант: «Мамина книга» (принадлежность)", []gramOption{
		{"妈妈的书", true}, {"妈妈书的", false}, {"的妈妈书", false}, {"书的妈妈", false},
	}, []string{"妈妈", "的", "书"}, "Притяжательная частица 的"},

	{"Как правильно сказать «У меня нет денег»?", []gramOption{
		{"我没有钱", true}, {"我不有钱", false}, {"我不是有钱", false}, {"我没是钱", false},
	}, []string{"我", "没", "有", "钱"}, "Отрицание: 不 и 没"},

	{"Выберите правильный вариант: «Это не большое»", []gramOption{
		{"这不大", true}, {"这没大", false}, {"不这大", false}, {"这不是大是", false},
	}, []string{"这", "不", "大"}, "Отрицание: 不 и 没"},

	{"Как правильно сказать «У меня есть один друг»?", []gramOption{
		{"我有一个朋友", true}, {"我有一本朋友", false}, {"我有朋友一个", false}, {"我有个一朋友", false},
	}, []string{"我", "有", "一", "个", "朋友"}, "Числительное + счётное слово + существительное"},

	{"Выберите правильный порядок слов: «Я сегодня иду в школу»", []gramOption{
		{"我今天去学校", true}, {"我去今天学校", false}, {"今天我学校去", false}, {"我去学校今天", false},
	}, []string{"我", "今天", "去", "学校"}, "Порядок слов: подлежащее + время + сказуемое + дополнение"},

	{"Как правильно сказать «Я хочу пить воду»?", []gramOption{
		{"我想喝水", true}, {"我喝想水", false}, {"想我喝水", false}, {"我喝水想", false},
	}, []string{"我", "想", "喝", "水"}, "Модальные глаголы перед основным глаголом"},

	{"Выберите правильный вариант: «Я тоже врач»", []gramOption{
		{"我也是医生", true}, {"我是也医生", false}, {"也我是医生", false}, {"我是医生也", false},
	}, []string{"我", "也", "是", "医生"}, "Наречия 也 и 都 перед сказуемым"},

	{"Выберите правильный вариант: «Мы все ученики»", []gramOption{
		{"我们都是学生", true}, {"都我们是学生", false}, {"我们是都学生", false}, {"我们是学生都", false},
	}, []string{"我们", "都", "是", "学生"}, "Наречия 也 и 都 перед сказуемым"},

	{"Как коротко спросить «А ты?» в ответ на «Я иду в школу»?", []gramOption{
		{"你呢?", true}, {"你吗?", false}, {"你了?", false}, {"你的?", false},
	}, []string{"你", "呢"}, "Краткий вопрос с 呢"},

	{"Как правильно сказать «Я уже поел» (завершённое действие)?", []gramOption{
		{"我吃了", true}, {"我了吃", false}, {"我吃的", false}, {"我吃吗", false},
	}, []string{"我", "吃", "了"}, "Частица 了 — завершённое действие"},

	{"Выберите правильный порядок слов: «Что ты ешь?»", []gramOption{
		{"你吃什么?", true}, {"什么你吃?", false}, {"你什么吃?", false}, {"吃你什么?", false},
	}, []string{"你", "吃", "什么"}, "Вопросительные слова: 什么, 谁, 哪儿"},

	{"Выберите правильный порядок слов: «Кто твой учитель?»", []gramOption{
		{"谁是你的老师?", true}, {"你的老师谁是?", false}, {"是谁你的老师?", false}, {"你谁的老师是?", false},
	}, []string{"谁", "是", "你", "的", "老师"}, "Притяжательная частица 的"},

	{"Как правильно сказать «Погода очень хорошая»?", []gramOption{
		{"天气很好", true}, {"天气好很", false}, {"很天气好", false}, {"天气很是好", false},
	}, []string{"天气", "很", "好"}, "很 перед прилагательным"},

	{"Как правильно сказать «Я нахожусь дома»?", []gramOption{
		{"我在家", true}, {"我家在", false}, {"在我家", false}, {"我是在家", false},
	}, []string{"我", "在", "家"}, "在 — местонахождение"},

	{"Выберите правильный порядок слов: «Где твой дом?»", []gramOption{
		{"你家在哪儿?", true}, {"哪儿你家在?", false}, {"你在哪儿家?", false}, {"你家哪儿在?", false},
	}, []string{"你", "家", "在", "哪儿"}, "在 — местонахождение"},

	{"Как правильно спросить цену: «Сколько это стоит?»", []gramOption{
		{"这个多少钱?", true}, {"这个几钱?", false}, {"多少这个钱?", false}, {"这个钱多少?", false},
	}, []string{"这", "个", "多少", "钱"}, "Вопросительные слова: 什么, 谁, 哪儿"},

	{"Как правильно спросить «Сколько у тебя книг?» (малое число + счётное слово)", []gramOption{
		{"你有几本书?", true}, {"你有几书?", false}, {"你有本几书?", false}, {"几你有本书?", false},
	}, []string{"你", "有", "几", "本", "书"}, "Числительное + счётное слово + существительное"},

	{"Как правильно сказать «Я умею говорить по-китайски»?", []gramOption{
		{"我会说汉语", true}, {"我说会汉语", false}, {"会我说汉语", false}, {"我会汉语说", false},
	}, []string{"我", "会", "说", "汉语"}, "Модальные глаголы перед основным глаголом"},
}

func main() {
	ctx := context.Background()

	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	var courseID int64

	if err := db.QueryRow(ctx, `SELECT id FROM courses WHERE title = 'HSK 1'`).Scan(&courseID); err != nil {
		log.Fatalf("failed to find HSK 1 course: %v", err)
	}

	if _, err := db.Exec(ctx, `
		DELETE FROM lessons WHERE course_id = $1 AND lesson_number >= $2
	`, courseID, firstGeneratedLesson); err != nil {
		log.Fatalf("failed to clear previously generated lessons: %v", err)
	}

	allWords, err := loadWords(ctx, db)
	if err != nil {
		log.Fatalf("failed to load words: %v", err)
	}

	byHanzi := make(map[string]wordInfo, len(allWords))
	translationByID := make(map[int64]string, len(allWords))
	allWordIDs := make([]int64, 0, len(allWords))

	for _, w := range allWords {
		byHanzi[w.Hanzi] = w
		translationByID[w.ID] = w.Translation
		allWordIDs = append(allWordIDs, w.ID)
	}

	introducedSet, err := loadIntroducedWordIDs(ctx, db, courseID)
	if err != nil {
		log.Fatalf("failed to load already-introduced words: %v", err)
	}

	introduced := make([]wordInfo, 0, len(allWords))
	for _, w := range allWords {
		if introducedSet[w.ID] {
			introduced = append(introduced, w)
		}
	}

	rng := rand.New(rand.NewSource(42))

	vocabQuizCache := map[int64]int64{}
	grammarQuizIDs := make([]int64, len(grammarBank))
	grammarCursor := 0
	shownConstructs := map[string]bool{}

	lessonNum := firstGeneratedLesson

	createOne := func(title, description string, newWords []wordInfo) {

		reviewPool := make([]wordInfo, 0, len(introduced))
		newIDs := make(map[int64]bool, len(newWords))
		for _, w := range newWords {
			newIDs[w.ID] = true
		}
		for _, w := range introduced {
			if !newIDs[w.ID] {
				reviewPool = append(reviewPool, w)
			}
		}

		rng.Shuffle(len(reviewPool), func(i, j int) {
			reviewPool[i], reviewPool[j] = reviewPool[j], reviewPool[i]
		})

		reviewCount := wordStepsPerLesson - len(newWords)
		if reviewCount > len(reviewPool) {
			reviewCount = len(reviewPool)
		}
		if reviewCount < 0 {
			reviewCount = 0
		}

		lessonWords := append(append([]wordInfo{}, newWords...), reviewPool[:reviewCount]...)

		// Pure-review lessons ("Итоговое повторение") introduce no new
		// words — guarantee quizzes for their whole review selection
		// instead, so they stay just as practice-heavy as before.
		guaranteed, bonus := newWords, bonusReviewQuizzes
		if len(newWords) == 0 {
			guaranteed, bonus = lessonWords, 0
		}

		// Everything a learner will have seen by the time they reach this
		// lesson's quiz steps: every previously introduced word, plus this
		// lesson's own new words (taught earlier in the same lesson).
		introducedHanzi := make(map[string]bool, len(introduced)+len(newWords))
		for _, w := range introduced {
			introducedHanzi[w.Hanzi] = true
		}
		for _, w := range newWords {
			introducedHanzi[w.Hanzi] = true
		}

		grammarEntries := selectGrammarEntries(&grammarCursor, grammarQuizzesPerLesson, introducedHanzi)

		if err := createLesson(ctx, db, courseID, lessonNum, title, description, lessonWords,
			guaranteed, bonus, grammarEntries, grammarQuizIDs, shownConstructs,
			vocabQuizCache, translationByID, allWordIDs, rng); err != nil {
			log.Fatalf("failed to create lesson %d (%s): %v", lessonNum, title, err)
		}

		fmt.Printf("Created lesson %d: %s (%d new, %d review, %d grammar)\n",
			lessonNum, title, len(newWords), reviewCount, len(grammarEntries))

		introduced = append(introduced, newWords...)
		lessonNum++
	}

	for _, topic := range topics {

		words := make([]wordInfo, 0, len(topic.hanzi))
		for _, h := range topic.hanzi {
			w, ok := byHanzi[h]
			if !ok {
				log.Fatalf("topic %q references unknown word %q", topic.name, h)
			}
			words = append(words, w)
		}

		part := 1
		for len(words) > 0 {

			n := 3
			if n > len(words) {
				n = len(words)
			}

			title := topic.name
			if part > 1 {
				title = fmt.Sprintf("%s %d", topic.name, part)
			}

			createOne(title, fmt.Sprintf("Новые слова: %s", topic.name), words[:n])

			words = words[n:]
			part++

			if lessonNum > totalLessons {
				break
			}
		}

		if lessonNum > totalLessons {
			break
		}
	}

	reviewPart := 1
	for lessonNum <= totalLessons {
		title := fmt.Sprintf("Итоговое повторение %d", reviewPart)
		createOne(title, "Повторение пройденных слов", nil)
		reviewPart++
	}

	// Grammar-explanation lesson + its sentence-builder review lesson (seeded
	// once by migration 000026/000029) are appended right after the
	// generated content, at whatever lesson_number comes next. Looked up by
	// id rather than duplicated, since the underlying grammar_notes/
	// sentence_exercises rows already exist.
	if err := createGrammarLesson(ctx, db, courseID, lessonNum); err != nil {
		log.Fatalf("failed to create grammar lesson %d: %v", lessonNum, err)
	}
	fmt.Printf("Created lesson %d: Как задавать вопросы\n", lessonNum)
	lessonNum++

	if err := createSentenceReviewLesson(ctx, db, courseID, lessonNum); err != nil {
		log.Fatalf("failed to create sentence-builder review lesson %d: %v", lessonNum, err)
	}
	fmt.Printf("Created lesson %d: Повторение: вопросы\n", lessonNum)

	fmt.Println("Done.")
}

// grammarUnlocked reports whether every word gq's example sentence relies on
// has already been introduced.
func grammarUnlocked(gq gramQuestion, introducedHanzi map[string]bool) bool {
	for _, h := range gq.requiredHanzi {
		if !introducedHanzi[h] {
			return false
		}
	}
	return true
}

// selectGrammarEntries scans grammarBank starting from *cursor (which
// persists across lessons, so it resumes where the previous lesson left
// off), returning up to `quota` indices whose required vocabulary is already
// introduced. Locked entries are skipped without being lost — the cursor
// simply passes over them and reconsiders them once it wraps back around,
// by which point more vocabulary is usually available. Scans at most one
// full pass over the bank per call, so a lesson early in the course (before
// much vocabulary exists) can legitimately return fewer than `quota`, or
// none at all.
func selectGrammarEntries(cursor *int, quota int, introducedHanzi map[string]bool) []int {

	n := len(grammarBank)
	selected := make([]int, 0, quota)

	for scanned := 0; len(selected) < quota && scanned < n; scanned++ {

		idx := *cursor % n
		*cursor++

		if grammarUnlocked(grammarBank[idx], introducedHanzi) {
			selected = append(selected, idx)
		}
	}

	return selected
}

func createGrammarLesson(ctx context.Context, db interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}, courseID int64, lessonNum int) error {

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var lessonID int64

	err = tx.QueryRow(ctx, `
		INSERT INTO lessons (course_id, title, description, lesson_number)
		VALUES ($1, 'Как задавать вопросы', 'Грамматика: вопрос с 吗 и вопросительные слова', $2)
		RETURNING id
	`, courseID, lessonNum).Scan(&lessonID)
	if err != nil {
		return err
	}

	rows, err := tx.Query(ctx, `SELECT id FROM grammar_notes ORDER BY id`)
	if err != nil {
		return err
	}

	noteIDs := make([]int64, 0)
	for rows.Next() {
		var noteID int64
		if err := rows.Scan(&noteID); err != nil {
			rows.Close()
			return err
		}
		noteIDs = append(noteIDs, noteID)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	for i, noteID := range noteIDs {
		if _, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, 'grammar', $2, $3)
		`, lessonID, noteID, i+1); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func createSentenceReviewLesson(ctx context.Context, db interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}, courseID int64, lessonNum int) error {

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var lessonID int64

	err = tx.QueryRow(ctx, `
		INSERT INTO lessons (course_id, title, description, lesson_number)
		VALUES ($1, 'Повторение: вопросы', 'Практика: собери вопросительное предложение', $2)
		RETURNING id
	`, courseID, lessonNum).Scan(&lessonID)
	if err != nil {
		return err
	}

	// The first 5 sentence_exercises rows are the original HSK1 question
	// sentences seeded by migration 000026; later migrations only append.
	rows, err := tx.Query(ctx, `SELECT id FROM sentence_exercises WHERE id BETWEEN 1 AND 5 ORDER BY id`)
	if err != nil {
		return err
	}

	exerciseIDs := make([]int64, 0)
	for rows.Next() {
		var exerciseID int64
		if err := rows.Scan(&exerciseID); err != nil {
			rows.Close()
			return err
		}
		exerciseIDs = append(exerciseIDs, exerciseID)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	for i, exerciseID := range exerciseIDs {
		if _, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, 'sentence_builder', $2, $3)
		`, lessonID, exerciseID, i+1); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func loadWords(ctx context.Context, db pgxPool) ([]wordInfo, error) {

	rows, err := db.Query(ctx, `SELECT id, hanzi, translation FROM words WHERE hsk_level = 1 ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	words := make([]wordInfo, 0)

	for rows.Next() {

		var w wordInfo

		if err := rows.Scan(&w.ID, &w.Hanzi, &w.Translation); err != nil {
			return nil, err
		}

		words = append(words, w)
	}

	return words, rows.Err()
}

func loadIntroducedWordIDs(ctx context.Context, db pgxPool, courseID int64) (map[int64]bool, error) {

	rows, err := db.Query(ctx, `
		SELECT DISTINCT ls.entity_id
		FROM lesson_steps ls
		JOIN lessons l ON l.id = ls.lesson_id
		WHERE l.course_id = $1
		  AND ls.step_type = 'word'
		  AND ls.entity_id IS NOT NULL
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	set := map[int64]bool{}

	for rows.Next() {

		var id int64

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		set[id] = true
	}

	return set, rows.Err()
}

func pickDistractors(rng *rand.Rand, allWordIDs []int64, translationByID map[int64]string, excludeID int64, n int) []string {

	shuffled := append([]int64{}, allWordIDs...)

	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	seen := map[string]bool{translationByID[excludeID]: true}
	result := make([]string, 0, n)

	for _, id := range shuffled {

		if id == excludeID {
			continue
		}

		t := translationByID[id]

		if seen[t] {
			continue
		}

		seen[t] = true
		result = append(result, t)

		if len(result) == n {
			break
		}
	}

	return result
}

// pgxPool is the subset of *pgxpool.Pool used by the loader helpers above.
type pgxPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func findQuizByQuestion(ctx context.Context, tx pgx.Tx, question string) (int64, bool, error) {

	var id int64

	err := tx.QueryRow(ctx, `SELECT id FROM quizzes WHERE question = $1`, question).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}

	return id, true, nil
}

// findGrammarNoteByTitle looks up a grammar_notes row seeded by a migration.
// Fatal if missing: every constructTitle referenced from grammarBank must
// have a matching row seeded ahead of time (migrations 000026, 000036).
func findGrammarNoteByTitle(ctx context.Context, tx pgx.Tx, title string) (int64, error) {

	var id int64

	err := tx.QueryRow(ctx, `SELECT id FROM grammar_notes WHERE title = $1`, title).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("grammar note %q not found (seed a migration for it first): %w", title, err)
	}

	return id, nil
}

func createLesson(
	ctx context.Context,
	db interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	},
	courseID int64,
	lessonNum int,
	title string,
	description string,
	lessonWords []wordInfo,
	newWords []wordInfo, // every one of these gets a quiz — guaranteed practice for what's actually taught this lesson
	bonusQuizCount int, // extra quizzes drawn from the remaining (review) words
	grammarEntries []int, // grammarBank indices already resolved (gated + selected) for this lesson
	grammarQuizIDs []int64,
	shownConstructs map[string]bool,
	vocabQuizCache map[int64]int64,
	translationByID map[int64]string,
	allWordIDs []int64,
	rng *rand.Rand,
) error {

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var lessonID int64

	err = tx.QueryRow(ctx, `
		INSERT INTO lessons (course_id, title, description, lesson_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, courseID, title, description, lessonNum).Scan(&lessonID)
	if err != nil {
		return err
	}

	sortOrder := 1

	for _, w := range lessonWords {

		_, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, 'word', $2, $3)
		`, lessonID, w.ID, sortOrder)
		if err != nil {
			return err
		}

		sortOrder++
	}

	newWordIDs := make(map[int64]bool, len(newWords))
	for _, w := range newWords {
		newWordIDs[w.ID] = true
	}

	reviewOnly := make([]wordInfo, 0, len(lessonWords))
	for _, w := range lessonWords {
		if !newWordIDs[w.ID] {
			reviewOnly = append(reviewOnly, w)
		}
	}

	rng.Shuffle(len(reviewOnly), func(i, j int) {
		reviewOnly[i], reviewOnly[j] = reviewOnly[j], reviewOnly[i]
	})

	if bonusQuizCount > len(reviewOnly) {
		bonusQuizCount = len(reviewOnly)
	}

	// Every new word is guaranteed a quiz; bonus review quizzes are appended
	// on top so practice never falls short of what was just taught.
	quizCandidates := append(append([]wordInfo{}, newWords...), reviewOnly[:bonusQuizCount]...)

	for i := 0; i < len(quizCandidates); i++ {

		w := quizCandidates[i]

		quizID, ok := vocabQuizCache[w.ID]

		if !ok {

			question := fmt.Sprintf("Как переводится %s?", w.Hanzi)

			existingID, found, err := findQuizByQuestion(ctx, tx, question)
			if err != nil {
				return err
			}

			if found {
				quizID = existingID
			} else {

				distractors := pickDistractors(rng, allWordIDs, translationByID, w.ID, 3)

				if err := tx.QueryRow(ctx, `
					INSERT INTO quizzes (question, hsk_level) VALUES ($1, 1) RETURNING id
				`, question).Scan(&quizID); err != nil {
					return err
				}

				options := append([]string{w.Translation}, distractors...)

				rng.Shuffle(len(options), func(a, b int) {
					options[a], options[b] = options[b], options[a]
				})

				for idx, optText := range options {

					_, err := tx.Exec(ctx, `
						INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
						VALUES ($1, $2, $3, $4)
					`, quizID, optText, optText == w.Translation, idx+1)
					if err != nil {
						return err
					}
				}
			}

			vocabQuizCache[w.ID] = quizID
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, 'quiz', $2, $3)
		`, lessonID, quizID, sortOrder)
		if err != nil {
			return err
		}

		sortOrder++
	}

	for _, idx := range grammarEntries {

		gq := grammarBank[idx]

		quizID := grammarQuizIDs[idx]

		if quizID == 0 {

			existingID, found, err := findQuizByQuestion(ctx, tx, gq.question)
			if err != nil {
				return err
			}

			if found {
				quizID = existingID
			} else {

				if err := tx.QueryRow(ctx, `
					INSERT INTO quizzes (question, hsk_level) VALUES ($1, 1) RETURNING id
				`, gq.question).Scan(&quizID); err != nil {
					return err
				}

				for optIdx, opt := range gq.options {

					_, err := tx.Exec(ctx, `
						INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
						VALUES ($1, $2, $3, $4)
					`, quizID, opt.text, opt.correct, optIdx+1)
					if err != nil {
						return err
					}
				}
			}

			grammarQuizIDs[idx] = quizID
		}

		// First time this construct is actually used: explain it before
		// testing it.
		if !shownConstructs[gq.constructTitle] {

			noteID, err := findGrammarNoteByTitle(ctx, tx, gq.constructTitle)
			if err != nil {
				return err
			}

			if _, err := tx.Exec(ctx, `
				INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
				VALUES ($1, 'grammar', $2, $3)
			`, lessonID, noteID, sortOrder); err != nil {
				return err
			}

			sortOrder++
			shownConstructs[gq.constructTitle] = true
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, 'quiz', $2, $3)
		`, lessonID, quizID, sortOrder)
		if err != nil {
			return err
		}

		sortOrder++
	}

	return tx.Commit(ctx)
}
