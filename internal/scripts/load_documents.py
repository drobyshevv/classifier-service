import sqlite3
import pandas as pd
from config import *

def load_documents():
    
    print("uploading documents")
    
    df = pd.read_csv(DOCUMENTS_CSV)
    
    conn = sqlite3.connect(DB_PATH)
    
    try:
        records = []
        for _, row in df.iterrows():
            record = (
                row['document_id'],
                row.get('linkurl'),
                row.get('genre'),
                row.get('type'),
                row.get('yearpubl'),
                row.get('pages'),
                row.get('language'),
                row.get('cited'),
                row.get('title_ru'),
                row.get('title_en'),
                row.get('edn'),
                row.get('grnti'),
                row.get('risc'),
                row.get('corerisc'),
                row.get('abstract_ru'),
                row.get('abstract_en'),
                row.get('codetype'),
                row.get('code'),
                row.get('reference'),
                row.get('numauthors'),
                row.get('title_j'),
                row.get('issn_j'),
                row.get('eissn_j'),
                row.get('publisher_j'),
                row.get('country_j'),
                row.get('town_j'),
                row.get('vak_j'),
                row.get('rsci_j'),
                row.get('wos_j'),
                row.get('scopus_j'),
                row.get('volume_j'),
                row.get('number_j'),
                row.get('created_at')
            )
            records.append(record)
        
        insert_sql = """
        INSERT OR REPLACE INTO documents (
            document_id, linkurl, genre, type, yearpubl, pages, language, cited,
            title_ru, title_en, edn, grnti, risc, corerisc, abstract_ru, abstract_en,
            codetype, code, reference, numauthors, title_j, issn_j, eissn_j, publisher_j,
            country_j, town_j, vak_j, rsci_j, wos_j, scopus_j, volume_j, number_j, created_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """
        
        conn.executemany(insert_sql, records)
        conn.commit()
        
        print(f"{len(records)} documents uploaded")
        
    except Exception as e:
        print(f"Error when uploading documents: {e}")
        conn.rollback()
    finally:
        conn.close()

if __name__ == "__main__":
    load_documents()