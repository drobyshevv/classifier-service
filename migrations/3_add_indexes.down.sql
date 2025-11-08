-- Удаляем индексы в обратном порядке
DROP INDEX IF EXISTS idx_categories_active;
DROP INDEX IF EXISTS idx_user_classifications_allowed;
DROP INDEX IF EXISTS idx_user_classifications_category;
DROP INDEX IF EXISTS idx_user_classifications_created_at;
DROP INDEX IF EXISTS idx_user_classifications_user_id;
DROP INDEX IF EXISTS idx_document_authors_org_id;
DROP INDEX IF EXISTS idx_document_authors_author_id;
DROP INDEX IF EXISTS idx_document_authors_document_id;
DROP INDEX IF EXISTS idx_documents_language;
DROP INDEX IF EXISTS idx_documents_year;