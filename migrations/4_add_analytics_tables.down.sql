-- Удаляем в обратном порядке из-за foreign keys
DROP TABLE IF EXISTS department_topics;
DROP TABLE IF EXISTS topic_experts;
DROP TABLE IF EXISTS article_topics;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_department_topics_org;
DROP INDEX IF EXISTS idx_department_topics_topic;
DROP INDEX IF EXISTS idx_topic_experts_author;
DROP INDEX IF EXISTS idx_topic_experts_topic;
DROP INDEX IF EXISTS idx_article_topics_document;
DROP INDEX IF EXISTS idx_article_topics_topic;