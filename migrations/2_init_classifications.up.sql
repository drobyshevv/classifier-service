-- История классификаций пользователей
CREATE TABLE user_classifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,  -- Связь с users.id из auth service
    document_id INTEGER,       -- Связь с существующими документами
    input_text TEXT NOT NULL,  -- Текст для классификации
    category TEXT NOT NULL,    -- Основная категория
    confidence REAL NOT NULL,  -- Уверенность модели (0.0-1.0)
    is_allowed BOOLEAN NOT NULL, -- Разрешена ли публикация
    themes JSON,               -- Дополнительные темы/метки
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(document_id)
);

-- Категории/темы для классификации
CREATE TABLE classification_categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE
);

-- Статистика по пользователям
CREATE TABLE user_classification_stats (
    user_id INTEGER PRIMARY KEY,
    total_classifications INTEGER DEFAULT 0,
    allowed_count INTEGER DEFAULT 0,
    rejected_count INTEGER DEFAULT 0,
    last_classification_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);