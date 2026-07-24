CREATE TABLE grammar_notes (
    id BIGSERIAL PRIMARY KEY,

    title VARCHAR(200) NOT NULL,
    explanation TEXT NOT NULL,

    example_hanzi VARCHAR(255),
    example_pinyin VARCHAR(255),
    example_translation VARCHAR(255),

    hsk_level SMALLINT NOT NULL DEFAULT 1,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sentence_exercises (
    id BIGSERIAL PRIMARY KEY,

    translation VARCHAR(255) NOT NULL,
    chunks TEXT[] NOT NULL,

    hsk_level SMALLINT NOT NULL DEFAULT 1,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO grammar_notes (title, explanation, example_hanzi, example_pinyin, example_translation, hsk_level) VALUES
('Вопрос с частицей 吗',
 'Чтобы превратить обычное утверждение в вопрос «да/нет», достаточно добавить частицу 吗 (ma) в самый конец предложения. Порядок слов при этом не меняется — никакой перестановки, как в русском или английском, не требуется.',
 '你是学生吗？', 'nǐ shì xué shēng ma?', 'Ты студент?', 1),
('Вопросительные слова: 什么, 谁, 哪儿',
 'Вопросительные слова 什么 (что), 谁 (кто), 哪儿 (где) в китайском предложении встают ровно на то же место, где в ответе стоял бы искомый элемент — само предложение не перестраивается. Сравните: «你叫什么名字» дословно «ты зовёшься каким именем» — 什么 стоит там же, где в ответе было бы имя.',
 '你叫什么名字？', 'nǐ jiào shén me míng zi?', 'Как тебя зовут?', 1);

INSERT INTO sentence_exercises (translation, chunks, hsk_level) VALUES
('Ты студент?', ARRAY['你','是','学生','吗'], 1),
('Как тебя зовут?', ARRAY['你','叫','什么','名字'], 1),
('Где твой дом?', ARRAY['你','家','在','哪儿'], 1),
('Сколько тебе лет?', ARRAY['你','今年','多大'], 1),
('Что ты ешь?', ARRAY['你','吃','什么'], 1);

-- Append two new lessons after the end of the HSK1 course (lesson_number
-- 51/52): a grammar lesson on question formation, immediately followed by
-- its review lesson containing a sentence-builder step. They're placed
-- after 50 rather than inserted mid-sequence because lessons 7..50 are
-- regenerated wholesale by cmd/seed and would wipe anything placed there.
INSERT INTO lessons (course_id, title, description, lesson_number)
SELECT c.id, 'Как задавать вопросы', 'Грамматика: вопрос с 吗 и вопросительные слова', 51
FROM courses c WHERE c.title = 'HSK 1';

INSERT INTO lessons (course_id, title, description, lesson_number)
SELECT c.id, 'Повторение: вопросы', 'Практика: собери вопросительное предложение', 52
FROM courses c WHERE c.title = 'HSK 1';

INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'grammar', gn.id, 1
FROM lessons l, grammar_notes gn
WHERE l.lesson_number = 51 AND l.course_id = (SELECT id FROM courses WHERE title = 'HSK 1')
  AND gn.title = 'Вопрос с частицей 吗';

INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'grammar', gn.id, 2
FROM lessons l, grammar_notes gn
WHERE l.lesson_number = 51 AND l.course_id = (SELECT id FROM courses WHERE title = 'HSK 1')
  AND gn.title = 'Вопросительные слова: 什么, 谁, 哪儿';

INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'sentence_builder', se.id, ROW_NUMBER() OVER (ORDER BY se.id)
FROM lessons l, sentence_exercises se
WHERE l.lesson_number = 52 AND l.course_id = (SELECT id FROM courses WHERE title = 'HSK 1');
