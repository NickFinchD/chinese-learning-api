CREATE TABLE saved_words (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    word_id BIGINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_saved_words_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_saved_words_word
        FOREIGN KEY (word_id)
        REFERENCES words(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_saved_word
        UNIQUE (user_id, word_id)
);