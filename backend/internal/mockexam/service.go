package mockexam

import (
	"context"
	"errors"
	"math/rand"

	"github.com/NickFinchD/chinese-learning-api/internal/quizzes"
	"github.com/NickFinchD/chinese-learning-api/internal/sentences"
)

// levelConfig sizes the simplified practice paper for a given HSK level.
// Counts are a proxy for the real exam's item counts, scaled down to what
// this app's content actually supports (no listening section — there's no
// audio question bank yet, only text-based vocab/grammar quizzes and
// sentence-builder exercises).
type levelConfig struct {
	quizCount     int
	sentenceCount int
	timeLimit     int // seconds
}

var levelConfigs = map[int16]levelConfig{
	1: {quizCount: 15, sentenceCount: 5, timeLimit: 15 * 60},
	2: {quizCount: 20, sentenceCount: 8, timeLimit: 25 * 60},
}

var errUnsupportedLevel = errors.New("mock exam not available for this HSK level yet")

type repository interface {
	CreateAttempt(ctx context.Context, userID int64, hskLevel int16, totalQuestions, correctCount, scorePercent int16, passed bool, durationSeconds int32) (*Attempt, error)
	HasPassed(ctx context.Context, userID int64, hskLevel int16) (bool, error)
	ListAttempts(ctx context.Context, userID int64) ([]Attempt, error)
}

type Service struct {
	repository       repository
	quizProvider     QuizProvider
	sentenceProvider SentenceProvider
	xpAwarder        XPAwarder
}

func NewService(
	repository repository,
	quizProvider QuizProvider,
	sentenceProvider SentenceProvider,
	xpAwarder XPAwarder,
) *Service {
	return &Service{
		repository:       repository,
		quizProvider:     quizProvider,
		sentenceProvider: sentenceProvider,
		xpAwarder:        xpAwarder,
	}
}

func shuffleQuizzes(items []quizzes.Quiz) {
	rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
}

func shuffleExercises(items []sentences.Exercise) {
	rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
}

func (s *Service) BuildPaper(ctx context.Context, hskLevel int16) (*Paper, error) {

	cfg, ok := levelConfigs[hskLevel]
	if !ok {
		return nil, errUnsupportedLevel
	}

	allQuizzes, err := s.quizProvider.GetByHSKLevel(ctx, hskLevel)
	if err != nil {
		return nil, err
	}

	allSentences, err := s.sentenceProvider.List(ctx, hskLevel)
	if err != nil {
		return nil, err
	}

	shuffleQuizzes(allQuizzes)
	shuffleExercises(allSentences)

	quizN := cfg.quizCount
	if quizN > len(allQuizzes) {
		quizN = len(allQuizzes)
	}

	sentenceN := cfg.sentenceCount
	if sentenceN > len(allSentences) {
		sentenceN = len(allSentences)
	}

	questions := make([]Question, 0, quizN+sentenceN)

	for i := 0; i < quizN; i++ {
		quiz := allQuizzes[i]
		questions = append(questions, Question{Type: "quiz", Quiz: &quiz})
	}

	for i := 0; i < sentenceN; i++ {
		exercise := allSentences[i]
		questions = append(questions, Question{Type: "sentence", Sentence: &exercise})
	}

	rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })

	return &Paper{
		HSKLevel:         hskLevel,
		TimeLimitSeconds: cfg.timeLimit,
		Questions:        questions,
	}, nil
}

func (s *Service) SubmitAttempt(
	ctx context.Context,
	userID int64,
	hskLevel int16,
	quizAnswers []QuizAnswer,
	sentenceAnswers []SentenceAnswer,
	durationSeconds int32,
) (*Result, error) {

	if _, ok := levelConfigs[hskLevel]; !ok {
		return nil, errUnsupportedLevel
	}

	total := len(quizAnswers) + len(sentenceAnswers)
	correct := 0

	for _, qa := range quizAnswers {

		ok, err := s.quizProvider.CheckAnswer(ctx, qa.QuizID, qa.OptionID)
		if err != nil {
			return nil, err
		}

		if ok {
			correct++
		}
	}

	for _, sa := range sentenceAnswers {

		ok, err := s.sentenceProvider.CheckAnswer(ctx, sa.ExerciseID, sa.Chunks)
		if err != nil {
			return nil, err
		}

		if ok {
			correct++
		}
	}

	scorePercent := 0
	if total > 0 {
		scorePercent = correct * 100 / total
	}

	passed := scorePercent >= PassThresholdPercent

	alreadyPassed, err := s.repository.HasPassed(ctx, userID, hskLevel)
	if err != nil {
		return nil, err
	}

	attempt, err := s.repository.CreateAttempt(
		ctx, userID, hskLevel,
		int16(total), int16(correct), int16(scorePercent),
		passed, durationSeconds,
	)
	if err != nil {
		return nil, err
	}

	xpAwarded := 0

	if passed && !alreadyPassed {
		if err := s.xpAwarder.AwardXP(ctx, userID, ExamPassXP); err != nil {
			return nil, err
		}
		xpAwarded = ExamPassXP
	}

	return &Result{Attempt: *attempt, XPAwarded: xpAwarded}, nil
}

func (s *Service) ListAttempts(ctx context.Context, userID int64) ([]Attempt, error) {
	return s.repository.ListAttempts(ctx, userID)
}
