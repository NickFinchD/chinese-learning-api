CREATE TABLE word_examples (
    id BIGSERIAL PRIMARY KEY,

    word_id BIGINT NOT NULL REFERENCES words(id) ON DELETE CASCADE,

    hanzi TEXT NOT NULL,
    pinyin TEXT NOT NULL,
    translation TEXT NOT NULL,

    sort_order SMALLINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_word_examples_word_id ON word_examples(word_id);
