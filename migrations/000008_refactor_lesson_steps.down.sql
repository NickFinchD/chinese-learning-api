ALTER TABLE lesson_steps
DROP COLUMN entity_id;

ALTER TABLE lesson_steps
ADD COLUMN word_id BIGINT,
ADD COLUMN grammar_rule_id BIGINT,
ADD COLUMN test_id BIGINT,
ADD COLUMN text_id BIGINT;