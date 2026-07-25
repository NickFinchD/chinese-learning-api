ALTER TABLE quizzes ADD COLUMN pinyin TEXT;

-- Auto-generated vocab quizzes always follow "Как переводится <hanzi>?" —
-- backfill pinyin for existing rows by matching that hanzi back to `words`.
UPDATE quizzes q
SET pinyin = w.pinyin
FROM words w
WHERE q.question = 'Как переводится ' || w.hanzi || '?'
  AND q.pinyin IS NULL;
