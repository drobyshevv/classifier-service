import sqlite3
import pandas as pd
from config import *

def load_authors():
    
    print("uploading authors")
    
    df = pd.read_csv(AUTHORS_CSV)
    
    conn = sqlite3.connect(DB_PATH)
    
    try:
        records = []
        for _, row in df.iterrows():
            record = (
                row['author_id'],
                row.get('lastname_ru'),
                row.get('initials_ru'),
                row.get('lastname_en'),
                row.get('initials_en'),
                bool(row.get('is_real_id', True)),
                row.get('caf')
            )
            records.append(record)
        
        insert_sql = """
        INSERT OR REPLACE INTO authors (
            author_id, lastname_ru, initials_ru, lastname_en, initials_en, is_real_id, caf
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
        """
        
        conn.executemany(insert_sql, records)
        conn.commit()
        
        print(f"{len(records)} authors uploaded")
        
    except Exception as e:
        print(f"error when uploading authors: {e}")
        conn.rollback()
    finally:
        conn.close()

if __name__ == "__main__":
    load_authors()