CREATE TABLE word_learning_progress (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    word_id BIGINT NOT NULL,

    stage SMALLINT NOT NULL DEFAULT 0,

    next_eligible_at TIMESTAMP,
    learned_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_word_learning_progress_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_word_learning_progress_word
        FOREIGN KEY (word_id)
        REFERENCES words(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_word_learning_progress_user_word
        UNIQUE (user_id, word_id)
);
