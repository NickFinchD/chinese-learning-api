-- Ensure the HSK1 course exists (it may already have been created manually
-- before migrations tracked seed data — reuse it instead of duplicating).
INSERT INTO courses (title, description, hsk_level, sort_order)
SELECT 'HSK 1', 'Начальный уровень', 1, 1
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'HSK 1');

-- Ensure all six HSK1 lessons exist for that course.
INSERT INTO lessons (course_id, title, description, lesson_number)
SELECT c.id, v.title, v.description, v.lesson_number
FROM courses c
CROSS JOIN (VALUES
    (1, 'Приветствие', 'Базовые фразы приветствия и прощания'),
    (2, 'Числа', 'Числа от одного до десяти'),
    (3, 'Семья', 'Члены семьи и личные местоимения'),
    (4, 'Глаголы и действия', 'Основные глаголы повседневной жизни'),
    (5, 'Время и даты', 'Дни, месяцы и время суток'),
    (6, 'Еда и напитки', 'Базовая еда и напитки')
) AS v(lesson_number, title, description)
WHERE c.title = 'HSK 1'
  AND NOT EXISTS (
    SELECT 1 FROM lessons l WHERE l.course_id = c.id AND l.lesson_number = v.lesson_number
  );

-- Lesson 2 "Числа": word steps 一..十
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'word', w.id, ws.ord
FROM lessons l
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
JOIN (VALUES
    ('一', 1), ('二', 2), ('三', 3), ('四', 4), ('五', 5),
    ('六', 6), ('七', 7), ('八', 8), ('九', 9), ('十', 10)
) AS ws(hanzi, ord) ON true
JOIN words w ON w.hanzi = ws.hanzi
WHERE l.lesson_number = 2;

WITH new_quiz AS (
    INSERT INTO quizzes (question, hsk_level)
    VALUES ('Как переводится 五?', 1)
    RETURNING id
), quiz_opts AS (
    INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
    SELECT new_quiz.id, o.text, o.correct, o.ord
    FROM new_quiz
    CROSS JOIN (VALUES
        ('Пять', TRUE, 1),
        ('Четыре', FALSE, 2),
        ('Шесть', FALSE, 3),
        ('Семь', FALSE, 4)
    ) AS o(text, correct, ord)
    RETURNING quiz_id
)
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'quiz', new_quiz.id, 11
FROM new_quiz
JOIN lessons l ON true
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
WHERE l.lesson_number = 2;

-- Lesson 3 "Семья": word steps for family members and pronouns
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'word', w.id, ws.ord
FROM lessons l
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
JOIN (VALUES
    ('爸爸', 1), ('妈妈', 2), ('儿子', 3), ('女儿', 4), ('我们', 5), ('你们', 6)
) AS ws(hanzi, ord) ON true
JOIN words w ON w.hanzi = ws.hanzi
WHERE l.lesson_number = 3;

WITH new_quiz AS (
    INSERT INTO quizzes (question, hsk_level)
    VALUES ('Как переводится 爸爸?', 1)
    RETURNING id
), quiz_opts AS (
    INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
    SELECT new_quiz.id, o.text, o.correct, o.ord
    FROM new_quiz
    CROSS JOIN (VALUES
        ('Папа', TRUE, 1),
        ('Мама', FALSE, 2),
        ('Сын', FALSE, 3),
        ('Дочь', FALSE, 4)
    ) AS o(text, correct, ord)
    RETURNING quiz_id
)
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'quiz', new_quiz.id, 7
FROM new_quiz
JOIN lessons l ON true
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
WHERE l.lesson_number = 3;

-- Lesson 4 "Глаголы и действия"
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'word', w.id, ws.ord
FROM lessons l
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
JOIN (VALUES
    ('吃', 1), ('喝', 2), ('看', 3), ('听', 4), ('说', 5), ('读', 6),
    ('写', 7), ('买', 8), ('去', 9), ('来', 10), ('做', 11), ('喜欢', 12)
) AS ws(hanzi, ord) ON true
JOIN words w ON w.hanzi = ws.hanzi
WHERE l.lesson_number = 4;

WITH new_quiz AS (
    INSERT INTO quizzes (question, hsk_level)
    VALUES ('Как переводится 喜欢?', 1)
    RETURNING id
), quiz_opts AS (
    INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
    SELECT new_quiz.id, o.text, o.correct, o.ord
    FROM new_quiz
    CROSS JOIN (VALUES
        ('Нравиться', TRUE, 1),
        ('Хотеть', FALSE, 2),
        ('Знать', FALSE, 3),
        ('Уметь', FALSE, 4)
    ) AS o(text, correct, ord)
    RETURNING quiz_id
)
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'quiz', new_quiz.id, 13
FROM new_quiz
JOIN lessons l ON true
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
WHERE l.lesson_number = 4;

-- Lesson 5 "Время и даты"
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'word', w.id, ws.ord
FROM lessons l
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
JOIN (VALUES
    ('今天', 1), ('明天', 2), ('昨天', 3), ('现在', 4),
    ('年', 5), ('月', 6), ('星期', 7), ('点', 8)
) AS ws(hanzi, ord) ON true
JOIN words w ON w.hanzi = ws.hanzi
WHERE l.lesson_number = 5;

WITH new_quiz AS (
    INSERT INTO quizzes (question, hsk_level)
    VALUES ('Как переводится 明天?', 1)
    RETURNING id
), quiz_opts AS (
    INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
    SELECT new_quiz.id, o.text, o.correct, o.ord
    FROM new_quiz
    CROSS JOIN (VALUES
        ('Завтра', TRUE, 1),
        ('Сегодня', FALSE, 2),
        ('Вчера', FALSE, 3),
        ('Сейчас', FALSE, 4)
    ) AS o(text, correct, ord)
    RETURNING quiz_id
)
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'quiz', new_quiz.id, 9
FROM new_quiz
JOIN lessons l ON true
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
WHERE l.lesson_number = 5;

-- Lesson 6 "Еда и напитки"
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'word', w.id, ws.ord
FROM lessons l
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
JOIN (VALUES
    ('茶', 1), ('水', 2), ('米饭', 3), ('菜', 4), ('水果', 5), ('苹果', 6)
) AS ws(hanzi, ord) ON true
JOIN words w ON w.hanzi = ws.hanzi
WHERE l.lesson_number = 6;

WITH new_quiz AS (
    INSERT INTO quizzes (question, hsk_level)
    VALUES ('Как переводится 米饭?', 1)
    RETURNING id
), quiz_opts AS (
    INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
    SELECT new_quiz.id, o.text, o.correct, o.ord
    FROM new_quiz
    CROSS JOIN (VALUES
        ('Варёный рис', TRUE, 1),
        ('Чай', FALSE, 2),
        ('Вода', FALSE, 3),
        ('Фрукты', FALSE, 4)
    ) AS o(text, correct, ord)
    RETURNING quiz_id
)
INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
SELECT l.id, 'quiz', new_quiz.id, 7
FROM new_quiz
JOIN lessons l ON true
JOIN courses c ON c.id = l.course_id AND c.title = 'HSK 1'
WHERE l.lesson_number = 6;
