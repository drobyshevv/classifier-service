import sqlite3
import pandas as pd
from config import *

def load_document_authors():
    
    print("uploading document-author links")
    
    df = pd.read_csv(DOCUMENT_AUTHORS_CSV)
    
    conn = sqlite3.connect(DB_PATH)
    
    try:
        records = []
        for _, row in df.iterrows():
            record = (
                row['id'],
                row['document_id'],
                row['author_id'],
                row['organization_id'],
                bool(row.get('first_author', False))
            )
            records.append(record)
        
        insert_sql = """
        INSERT OR REPLACE INTO document_authors (
            id, document_id, author_id, organization_id, first_author
        ) VALUES (?, ?, ?, ?, ?)
        """
        
        conn.executemany(insert_sql, records)
        conn.commit()
        
        print(f"{len(records)} links uploaded")
        
    except Exception as e:
        print(f"error loading links: {e}")
        conn.rollback()
    finally:
        conn.close()

if __name__ == "__main__":
    load_document_authors()