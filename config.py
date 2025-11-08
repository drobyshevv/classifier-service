import os

# Конфигурация
DB_PATH = 'classifier.db'
CSV_FOLDER = 'csv_data'  # Папка с CSV в корне проекта

# Пути к CSV файлам
DOCUMENTS_CSV = os.path.join(CSV_FOLDER, 'df_doc.csv')
AUTHORS_CSV = os.path.join(CSV_FOLDER, 'df_auth.csv')
ORGANIZATIONS_CSV = os.path.join(CSV_FOLDER, 'df_org.csv')
DOCUMENT_AUTHORS_CSV = os.path.join(CSV_FOLDER, 'df_doc_auth.csv')