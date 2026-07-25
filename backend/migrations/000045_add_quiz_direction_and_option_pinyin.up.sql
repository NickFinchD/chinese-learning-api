ALTER TABLE quizzes ADD COLUMN direction VARCHAR(30) NOT NULL DEFAULT 'word_to_translation';
ALTER TABLE quizzes ADD COLUMN hanzi TEXT;
ALTER TABLE quiz_options ADD COLUMN pinyin TEXT;

-- Backfill hanzi for existing "Как переводится <hanzi>?" quizzes, mirroring
-- the pinyin backfill in migration 000044.
UPDATE quizzes q
SET hanzi = w.hanzi
FROM words w
WHERE q.question = 'Как переводится ' || w.hanzi || '?'
  AND q.hanzi IS NULL;
