CREATE TABLE IF NOT EXISTS words (
    id SERIAL PRIMARY KEY,
    word TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_words (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    word_id INTEGER NOT NULL REFERENCES words(id),
    knowledge_level INTEGER NOT NULL CHECK (knowledge_level >= 0 AND knowledge_level <= 5),
    UNIQUE(user_id, word_id)
);