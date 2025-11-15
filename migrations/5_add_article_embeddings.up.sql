-- Таблица для векторных представлений статей для семантического поиска (инициализация при запуске приложения)
CREATE TABLE IF NOT EXISTS article_embeddings (
    document_id INTEGER PRIMARY KEY,
    title_embedding BLOB,        -- Векторное представление заголовка статьи
    abstract_embedding BLOB,     -- Векторное представление аннотации статьи
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(document_id)
);

-- Индексы для быстрого поиска по векторам
CREATE INDEX IF NOT EXISTS idx_article_embeddings_created ON article_embeddings(created_at);