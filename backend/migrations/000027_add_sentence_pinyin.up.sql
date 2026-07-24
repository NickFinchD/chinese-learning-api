ALTER TABLE sentence_exercises ADD COLUMN pinyin VARCHAR(255) NOT NULL DEFAULT '';

UPDATE sentence_exercises SET pinyin = 'nǐ shì xué shēng ma?' WHERE translation = 'Ты студент?';
UPDATE sentence_exercises SET pinyin = 'nǐ jiào shén me míng zi?' WHERE translation = 'Как тебя зовут?';
UPDATE sentence_exercises SET pinyin = 'nǐ jiā zài nǎr?' WHERE translation = 'Где твой дом?';
UPDATE sentence_exercises SET pinyin = 'nǐ jīn nián duō dà?' WHERE translation = 'Сколько тебе лет?';
UPDATE sentence_exercises SET pinyin = 'nǐ chī shén me?' WHERE translation = 'Что ты ешь?';

ALTER TABLE sentence_exercises ALTER COLUMN pinyin DROP DEFAULT;
