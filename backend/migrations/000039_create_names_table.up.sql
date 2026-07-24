-- Proper nouns (people's names, place names) that appear in reading texts.
-- Kept separate from `words` on purpose: they get a hover hint (pinyin +
-- translation) in the reader like any dictionary word, but aren't real
-- vocabulary to learn/train/quiz on, so they never appear in the words
-- list, saved words, collections, or quizzes.
CREATE TABLE names (
    id BIGSERIAL PRIMARY KEY,

    hanzi VARCHAR(100) NOT NULL,
    pinyin VARCHAR(255) NOT NULL,
    translation VARCHAR(255) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
