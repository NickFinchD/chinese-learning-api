CREATE TABLE mock_exam_attempts (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    hsk_level SMALLINT NOT NULL,

    total_questions SMALLINT NOT NULL,
    correct_count SMALLINT NOT NULL,
    score_percent SMALLINT NOT NULL,
    passed BOOLEAN NOT NULL,

    duration_seconds INTEGER NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_mock_exam_attempts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_mock_exam_attempts_user_level ON mock_exam_attempts(user_id, hsk_level);
