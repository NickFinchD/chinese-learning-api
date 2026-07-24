CREATE TABLE user_text_progress (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    text_id BIGINT NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'in_progress',

    read_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_text_progress_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_text_progress_text
        FOREIGN KEY (text_id)
        REFERENCES texts(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_user_text
        UNIQUE(user_id, text_id)
);
