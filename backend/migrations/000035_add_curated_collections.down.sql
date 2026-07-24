DELETE FROM word_collection_items
WHERE collection_id IN (SELECT id FROM word_collections WHERE is_curated = TRUE);

DELETE FROM word_collections WHERE is_curated = TRUE;

ALTER TABLE word_collections DROP COLUMN source_collection_id;
ALTER TABLE word_collections DROP COLUMN is_curated;
ALTER TABLE word_collections ALTER COLUMN user_id SET NOT NULL;
