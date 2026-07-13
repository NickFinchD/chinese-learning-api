CREATE TABLE user_lesson_progress (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    lesson_id BIGINT NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'not_started',

    current_step INTEGER NOT NULL DEFAULT 0,
    score INTEGER NOT NULL DEFAULT 0,

    started_at TIMESTAMP,
    completed_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_progress_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_progress_lesson
        FOREIGN KEY (lesson_id)
        REFERENCES lessons(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_user_lesson
        UNIQUE(user_id, lesson_id)
);