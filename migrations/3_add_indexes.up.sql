-- Индексы для документов
CREATE INDEX IF NOT EXISTS idx_documents_year ON documents(yearpubl);
CREATE INDEX IF NOT EXISTS idx_documents_language ON documents(language);

-- Индексы для связей авторов
CREATE INDEX IF NOT EXISTS idx_document_authors_document_id ON document_authors(document_id);
CREATE INDEX IF NOT EXISTS idx_document_authors_author_id ON document_authors(author_id);
CREATE INDEX IF NOT EXISTS idx_document_authors_org_id ON document_authors(organization_id);

-- Индексы для классификаций
CREATE INDEX IF NOT EXISTS idx_user_classifications_user_id ON user_classifications(user_id);
CREATE INDEX IF NOT EXISTS idx_user_classifications_created_at ON user_classifications(created_at);
CREATE INDEX IF NOT EXISTS idx_user_classifications_category ON user_classifications(category);
CREATE INDEX IF NOT EXISTS idx_user_classifications_allowed ON user_classifications(is_allowed);

-- Индекс для категорий
CREATE INDEX IF NOT EXISTS idx_categories_active ON classification_categories(is_active);