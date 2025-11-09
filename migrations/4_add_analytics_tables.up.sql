-- Тематики статей (определенные AI)
CREATE TABLE IF NOT EXISTS article_topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    topic_name TEXT NOT NULL,           -- "машинное обучение"
    confidence REAL NOT NULL,           -- Уверенность AI
    topic_type TEXT NOT NULL,           -- "main", "secondary"
    FOREIGN KEY (document_id) REFERENCES documents(document_id)
);

-- Эксперты по темам (авторы)
CREATE TABLE IF NOT EXISTS topic_experts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id INTEGER NOT NULL,
    topic_name TEXT NOT NULL,
    expertise_score REAL NOT NULL,      -- Оценка экспертизы (0.0-1.0)
    article_count INTEGER NOT NULL,     -- Количество статей по теме
    total_citations INTEGER NOT NULL,   -- Сумма цитирований
    last_activity_year INTEGER,         -- Год последней публикации
    FOREIGN KEY (author_id) REFERENCES authors(author_id),
    UNIQUE(author_id, topic_name)
);

-- Тематики кафедр/подразделений
CREATE TABLE IF NOT EXISTS department_topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    topic_name TEXT NOT NULL,
    strength_score REAL NOT NULL,       -- Сила тематики в подразделении
    expert_count INTEGER NOT NULL,      -- Количество экспертов
    total_articles INTEGER NOT NULL,    -- Всего статей по теме
    FOREIGN KEY (organization_id) REFERENCES organizations(organization_id),
    UNIQUE(organization_id, topic_name)
);

-- Индексы для поиска
CREATE INDEX IF NOT EXISTS idx_article_topics_topic ON article_topics(topic_name);
CREATE INDEX IF NOT EXISTS idx_article_topics_document ON article_topics(document_id);
CREATE INDEX IF NOT EXISTS idx_topic_experts_topic ON topic_experts(topic_name);
CREATE INDEX IF NOT EXISTS idx_topic_experts_author ON topic_experts(author_id);
CREATE INDEX IF NOT EXISTS idx_department_topics_topic ON department_topics(topic_name);
CREATE INDEX IF NOT EXISTS idx_department_topics_org ON department_topics(organization_id);