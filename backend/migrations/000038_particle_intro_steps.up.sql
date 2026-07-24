-- Support up to 3 examples per grammar note, so a dedicated "examples" note
-- can show several usage sentences in a single step (the original single
-- example_* columns keep working unchanged for existing "principle" notes).
ALTER TABLE grammar_notes ADD COLUMN example2_hanzi VARCHAR(255);
ALTER TABLE grammar_notes ADD COLUMN example2_pinyin VARCHAR(255);
ALTER TABLE grammar_notes ADD COLUMN example2_translation VARCHAR(255);
ALTER TABLE grammar_notes ADD COLUMN example3_hanzi VARCHAR(255);
ALTER TABLE grammar_notes ADD COLUMN example3_pinyin VARCHAR(255);
ALTER TABLE grammar_notes ADD COLUMN example3_translation VARCHAR(255);

-- A handful of existing notes were only ever shown once their grammar-quiz
-- became eligible (i.e. once a whole example sentence's vocabulary was
-- taught), so their example reached for later words (妈妈, 书, 学生, 吃...).
-- cmd/seed now also shows the matching "principle" note the moment the
-- particle/auxiliary verb itself is introduced, much earlier — so these
-- examples are replaced with ones built only from pronouns and other
-- particles taught in the same early batch.
UPDATE grammar_notes SET example_hanzi = '我不是你', example_pinyin = 'wǒ bú shì nǐ', example_translation = 'Я не ты'
WHERE title = 'Отрицание: 不 и 没';

UPDATE grammar_notes SET example_hanzi = '这是我的', example_pinyin = 'zhè shì wǒ de', example_translation = 'Это моё'
WHERE title = 'Притяжательная частица 的';

UPDATE grammar_notes SET example_hanzi = '我和你都不是他', example_pinyin = 'wǒ hé nǐ dōu bú shì tā', example_translation = 'Ни я, ни ты — не он'
WHERE title = 'Наречия 也 и 都 перед сказуемым';

UPDATE grammar_notes SET example_hanzi = '我有你了', example_pinyin = 'wǒ yǒu nǐ le', example_translation = 'Теперь у меня есть ты'
WHERE title = 'Частица 了 — завершённое действие';

UPDATE grammar_notes SET example_hanzi = '这是你吗？', example_pinyin = 'zhè shì nǐ ma?', example_translation = 'Это ты?'
WHERE title = 'Вопрос с частицей 吗';

UPDATE grammar_notes SET example_hanzi = '我会说', example_pinyin = 'wǒ huì shuō', example_translation = 'Я умею говорить'
WHERE title = 'Модальные глаголы перед основным глаголом';

-- Three particles/verbs from the "Частицы и грамматика" topic didn't have a
-- dedicated principle note yet (是 and 有 were only ever used inside other
-- notes' examples; 和 wasn't covered at all).
INSERT INTO grammar_notes (title, explanation, example_hanzi, example_pinyin, example_translation, hsk_level) VALUES
('Глагол-связка 是 («быть»)',
 'Глагол 是 (shì) связывает подлежащее с существительным или местоимением, показывая, что это одно и то же: «подлежащее + 是 + существительное». Перед прилагательными 是 не используется — там своя конструкция с 很 (см. соответствующую заметку).',
 '这是我', 'zhè shì wǒ', 'Это я', 1),

('Глагол 有 — обладание и наличие',
 'Глагол 有 (yǒu) означает «иметь» (что-то есть у кого-то) или «существовать, находиться» (где-то что-то есть). Отрицается только через 没 (没有) — никогда через 不.',
 '我有你', 'wǒ yǒu nǐ', 'У меня есть ты', 1),

('Союз 和 — соединение существительных',
 'Союз 和 (hé, «и») соединяет однородные существительные или местоимения: «А 和 Б». В отличие от русского, 和 обычно не используется для соединения целых предложений или глаголов — там части предложения просто идут подряд, без союза.',
 '我和你', 'wǒ hé nǐ', 'Я и ты', 1);

-- Dedicated "examples" notes: shown as the step right after each principle
-- note, using only vocabulary already taught by that point in the course
-- (pronouns, and whichever of 是/不/有/的/都/和/了/吗/没/呢 came earlier in
-- the same topic batch).
INSERT INTO grammar_notes (title, explanation, example_hanzi, example_pinyin, example_translation, example2_hanzi, example2_pinyin, example2_translation, hsk_level) VALUES
('Примеры: 是',
 'Ещё примеры:',
 '那是你', 'nà shì nǐ', 'То — ты',
 '这是他', 'zhè shì tā', 'Это он',
 1),

('Примеры: 不 и 没',
 'Ещё примеры:',
 '他不是我', 'tā bú shì wǒ', 'Он не я',
 '这不是你', 'zhè bú shì nǐ', 'Это не ты',
 1),

('Примеры: 有',
 'Ещё примеры:',
 '他有你', 'tā yǒu nǐ', 'У него есть ты',
 '你有他', 'nǐ yǒu tā', 'У тебя есть он',
 1),

('Примеры: 的',
 'Ещё примеры:',
 '那是你的', 'nà shì nǐ de', 'То — твоё',
 '这是他的', 'zhè shì tā de', 'Это его',
 1),

('Примеры: 也 и 都',
 'Ещё примеры:',
 '你和他都是我的', 'nǐ hé tā dōu shì wǒ de', 'Ты и он — оба мои',
 '这和那都是你的', 'zhè hé nà dōu shì nǐ de', 'Это и то — оба твои',
 1),

('Примеры: 和',
 'Ещё примеры:',
 '他和她', 'tā hé tā', 'Он и она',
 '这和那', 'zhè hé nà', 'Это и то',
 1),

('Примеры: 了',
 'Ещё примеры:',
 '他有我了', 'tā yǒu wǒ le', 'Теперь у него есть я',
 '你有他了', 'nǐ yǒu tā le', 'Теперь у тебя есть он',
 1),

('Примеры: 吗',
 'Ещё примеры:',
 '你有我吗？', 'nǐ yǒu wǒ ma?', 'У тебя есть я?',
 '他是你的吗？', 'tā shì nǐ de ma?', 'Он твой?',
 1),

('Примеры: 呢',
 'Ещё примеры:',
 '他呢？', 'tā ne?', 'А он?',
 '那呢？', 'nà ne?', 'А то?',
 1),

('Примеры: модальных глаголов',
 'Ещё примеры:',
 '你会做吗？', 'nǐ huì zuò ma?', 'Ты умеешь готовить?',
 '他会写', 'tā huì xiě', 'Он умеет писать',
 1);
