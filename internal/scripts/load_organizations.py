import sqlite3
import pandas as pd
from config import *

def load_organizations():
    
    print("uploading organizations")
    
    df = pd.read_csv(ORGANIZATIONS_CSV)
    
    conn = sqlite3.connect(DB_PATH)
    
    try:
        records = []
        for _, row in df.iterrows():
            record = (
                row['organization_id'],
                row.get('name'),
                bool(row.get('is_real_id', True))
            )
            records.append(record)
        
        insert_sql = """
        INSERT OR REPLACE INTO organizations (
            organization_id, name, is_real_id
        ) VALUES (?, ?, ?)
        """
        
        conn.executemany(insert_sql, records)
        conn.commit()
        
        print(f"{len(records)} organizations uploaded")
        
    except Exception as e:
        print(f"error when uploading organizations: {e}")
        conn.rollback()
    finally:
        conn.close()

if __name__ == "__main__":
    load_organizations()