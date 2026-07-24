-- Ensure the HSK2 course exists. Lessons are generated separately by
-- cmd/seed-hsk2 from the HSK2 words already seeded in migration 000020.
INSERT INTO courses (title, description, hsk_level, sort_order)
SELECT 'HSK 2', 'Средний уровень', 2, 2
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'HSK 2');
