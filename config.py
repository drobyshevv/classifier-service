import os

# Конфигурация
DB_PATH = '/Users/ilya/vsCode/classifier-expert-search/storage/service.db'
CSV_FOLDER = '/Users/ilya/vsCode/classifier-expert-search/csv_data'

# Пути к CSV файлам
DOCUMENTS_CSV = os.path.join(CSV_FOLDER, 'df_doc.csv')
AUTHORS_CSV = os.path.join(CSV_FOLDER, 'df_auth.csv')
ORGANIZATIONS_CSV = os.path.join(CSV_FOLDER, 'df_org.csv')
DOCUMENT_AUTHORS_CSV = os.path.join(CSV_FOLDER, 'df_doc_auth.csv')