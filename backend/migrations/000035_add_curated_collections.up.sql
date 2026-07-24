ALTER TABLE word_collections ALTER COLUMN user_id DROP NOT NULL;
ALTER TABLE word_collections ADD COLUMN is_curated BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE word_collections ADD COLUMN source_collection_id BIGINT
    REFERENCES word_collections(id) ON DELETE SET NULL;

-- Starter thematic collections, owned by nobody (user_id NULL), visible to
-- every user as suggestions on the vocabulary page. Users "save" one to get
-- their own editable copy (word_collections row with source_collection_id
-- pointing back here) rather than editing the shared original.
INSERT INTO word_collections (user_id, name, is_curated) VALUES
(NULL, 'Семья', TRUE),
(NULL, 'Еда и напитки', TRUE),
(NULL, 'Числа', TRUE),
(NULL, 'Дом и вещи', TRUE);

INSERT INTO word_collection_items (collection_id, word_id)
SELECT c.id, w.id
FROM word_collections c
JOIN (VALUES
    ('Семья', '爸爸'), ('Семья', '妈妈'), ('Семья', '儿子'), ('Семья', '女儿'), ('Семья', '我们'), ('Семья', '你们'),
    ('Еда и напитки', '茶'), ('Еда и напитки', '水'), ('Еда и напитки', '米饭'), ('Еда и напитки', '菜'), ('Еда и напитки', '水果'), ('Еда и напитки', '苹果'),
    ('Числа', '一'), ('Числа', '二'), ('Числа', '三'), ('Числа', '四'), ('Числа', '五'), ('Числа', '六'), ('Числа', '七'), ('Числа', '八'), ('Числа', '九'), ('Числа', '十'),
    ('Дом и вещи', '杯子'), ('Дом и вещи', '电脑'), ('Дом и вещи', '电视'), ('Дом и вещи', '电影'), ('Дом и вещи', '东西'), ('Дом и вещи', '家'), ('Дом и вещи', '书'), ('Дом и вещи', '衣服'), ('Дом и вещи', '椅子'), ('Дом и вещи', '桌子')
) AS pairs(collection_name, hanzi) ON pairs.collection_name = c.name
JOIN words w ON w.hanzi = pairs.hanzi AND w.hsk_level = 1
WHERE c.is_curated = TRUE;
