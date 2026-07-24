CREATE TABLE word_collections (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_word_collections_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_word_collections_user ON word_collections(user_id);

CREATE TABLE word_collection_items (
    id BIGSERIAL PRIMARY KEY,

    collection_id BIGINT NOT NULL,
    word_id BIGINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_collection_items_collection
        FOREIGN KEY (collection_id)
        REFERENCES word_collections(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_collection_items_word
        FOREIGN KEY (word_id)
        REFERENCES words(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_collection_word
        UNIQUE (collection_id, word_id)
);

CREATE INDEX idx_collection_items_collection ON word_collection_items(collection_id);
