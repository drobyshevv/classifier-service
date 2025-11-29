import sqlite3
import os

DB_PATH = '/Users/ilya/vsCode/classifier-expert-search/storage/service.db'

def reset_tables():
    if not os.path.exists(DB_PATH):
        print(f"Database not found: {DB_PATH}")
        return
    
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    try:
        print("Cleaning tables...")

        tables = [
            "article_embeddings",
            "article_topics",
            "topic_experts",
            "department_topics",
        ]

        for t in tables:
            cursor.execute(f"DELETE FROM {t};")
            print(f"✔ Cleared table: {t}")

        print("Resetting AUTOINCREMENT counters...")
        cursor.execute("DELETE FROM sqlite_sequence WHERE name IN ('article_embeddings', 'article_topics', 'topic_experts', 'department_topics');")

        conn.commit()

        print("🧹 Running VACUUM (this may take a few seconds)...")
        cursor.execute("VACUUM;")

        print("\nReset completed successfully!")

    except Exception as e:
        print(f"Error: {e}")

    finally:
        cursor.close()
        conn.close()


if __name__ == "__main__":
    reset_tables()
