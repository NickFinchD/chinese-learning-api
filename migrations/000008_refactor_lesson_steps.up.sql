ALTER TABLE lesson_steps
DROP COLUMN word_id,
DROP COLUMN grammar_rule_id,
DROP COLUMN test_id,
DROP COLUMN text_id;

ALTER TABLE lesson_steps
ADD COLUMN entity_id BIGINT;