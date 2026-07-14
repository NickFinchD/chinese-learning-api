CREATE TABLE lesson_steps (
    id BIGSERIAL PRIMARY KEY,

    lesson_id BIGINT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,

    step_type VARCHAR(50) NOT NULL,

    sort_order INT NOT NULL,

    word_id BIGINT,
    grammar_rule_id BIGINT,
    test_id BIGINT,
    text_id BIGINT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);