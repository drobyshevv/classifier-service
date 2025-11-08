-- Документы/публикации (из df_doc.csv)
CREATE TABLE documents (
    document_id INTEGER PRIMARY KEY,
    linkurl TEXT,
    genre TEXT,
    type TEXT,
    yearpubl INTEGER,
    pages TEXT,
    language TEXT,
    cited INTEGER,
    title_ru TEXT,
    title_en TEXT,
    edn TEXT,
    grnti TEXT,
    risc TEXT,
    corerisc TEXT,
    abstract_ru TEXT,
    abstract_en TEXT,
    codetype TEXT,
    code TEXT,
    reference TEXT,
    numauthors INTEGER,
    title_j TEXT,
    issn_j TEXT,
    eissn_j TEXT,
    publisher_j TEXT,
    country_j TEXT,
    town_j TEXT,
    vak_j TEXT,
    rsci_j TEXT,
    wos_j TEXT,
    scopus_j TEXT,
    volume_j TEXT,
    number_j TEXT,
    created_at DATETIME
);

-- Авторы (из df_auth.csv)
CREATE TABLE authors (
    author_id INTEGER PRIMARY KEY,
    lastname_ru TEXT,
    initials_ru TEXT,
    lastname_en TEXT,
    initials_en TEXT,
    is_real_id BOOLEAN,
    caf INTEGER
);

-- Организации (из df_org.csv)
CREATE TABLE organizations (
    organization_id INTEGER PRIMARY KEY,
    name TEXT,
    is_real_id BOOLEAN
);

-- Связь документов с авторами и организациями (из df_doc_auth.csv)
CREATE TABLE document_authors (
    id INTEGER PRIMARY KEY,
    document_id INTEGER,
    author_id INTEGER,
    organization_id INTEGER,
    first_author BOOLEAN,
    FOREIGN KEY (document_id) REFERENCES documents(document_id),
    FOREIGN KEY (author_id) REFERENCES authors(author_id),
    FOREIGN KEY (organization_id) REFERENCES organizations(organization_id)
);