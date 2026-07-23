CREATE TABLE texts (
    id BIGSERIAL PRIMARY KEY,

    title VARCHAR(255) NOT NULL,
    hanzi TEXT NOT NULL,
    pinyin TEXT NOT NULL,
    translation TEXT NOT NULL,

    hsk_level SMALLINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
