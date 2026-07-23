-- The review module (separate spaced-repetition flashcards on saved words)
-- has been removed in favor of the learning module's word_learning_progress
-- table, which now tracks word repetition through the "Тесты" section.
DROP TABLE IF EXISTS user_word_progress;
