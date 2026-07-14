CREATE TABLE user_word_progress (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    word_id BIGINT NOT NULL,

    review_count INTEGER NOT NULL DEFAULT 0,
    correct_count INTEGER NOT NULL DEFAULT 0,
    wrong_count INTEGER NOT NULL DEFAULT 0,

    ease_factor NUMERIC(3,2) NOT NULL DEFAULT 2.50,

    last_review_at TIMESTAMP,
    next_review_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_word_progress_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_word_progress_word
        FOREIGN KEY (word_id)
        REFERENCES words(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_user_word
        UNIQUE(user_id, word_id)
);