DELETE FROM lessons WHERE lesson_number IN (51, 52) AND course_id = (SELECT id FROM courses WHERE title = 'HSK 1');

DROP TABLE sentence_exercises;
DROP TABLE grammar_notes;
