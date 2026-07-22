CREATE TABLE quizzes (
    id BIGSERIAL PRIMARY KEY,

    question TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE quiz_options (
    id BIGSERIAL PRIMARY KEY,

    quiz_id BIGINT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,

    option_text TEXT NOT NULL,

    is_correct BOOLEAN NOT NULL DEFAULT FALSE,

    sort_order SMALLINT NOT NULL
);