DELETE FROM grammar_notes WHERE title IN (
'Примеры: 是', 'Примеры: 不 и 没', 'Примеры: 有', 'Примеры: 的',
'Примеры: 也 и 都', 'Примеры: 和', 'Примеры: 了', 'Примеры: 吗',
'Примеры: 呢', 'Примеры: модальных глаголов'
);

DELETE FROM grammar_notes WHERE title IN (
'Глагол-связка 是 («быть»)', 'Глагол 有 — обладание и наличие', 'Союз 和 — соединение существительных'
);

UPDATE grammar_notes SET example_hanzi = '这不大', example_pinyin = 'zhè bú dà', example_translation = 'Это не большое'
WHERE title = 'Отрицание: 不 и 没';

UPDATE grammar_notes SET example_hanzi = '妈妈的书', example_pinyin = 'māma de shū', example_translation = 'Мамина книга'
WHERE title = 'Притяжательная частица 的';

UPDATE grammar_notes SET example_hanzi = '我们都是学生', example_pinyin = 'wǒmen dōu shì xuésheng', example_translation = 'Мы все ученики'
WHERE title = 'Наречия 也 и 都 перед сказуемым';

UPDATE grammar_notes SET example_hanzi = '我吃了', example_pinyin = 'wǒ chī le', example_translation = 'Я уже поел'
WHERE title = 'Частица 了 — завершённое действие';

UPDATE grammar_notes SET example_hanzi = '你是学生吗？', example_pinyin = 'nǐ shì xué shēng ma?', example_translation = 'Ты студент?'
WHERE title = 'Вопрос с частицей 吗';

UPDATE grammar_notes SET example_hanzi = '我会说汉语', example_pinyin = 'wǒ huì shuō hànyǔ', example_translation = 'Я умею говорить по-китайски'
WHERE title = 'Модальные глаголы перед основным глаголом';

ALTER TABLE grammar_notes DROP COLUMN example3_translation;
ALTER TABLE grammar_notes DROP COLUMN example3_pinyin;
ALTER TABLE grammar_notes DROP COLUMN example3_hanzi;
ALTER TABLE grammar_notes DROP COLUMN example2_translation;
ALTER TABLE grammar_notes DROP COLUMN example2_pinyin;
ALTER TABLE grammar_notes DROP COLUMN example2_hanzi;
