import sys
import os

sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(__file__))))

from config import *
from load_documents import load_documents
from load_authors import load_authors
from load_organizations import load_organizations
from load_document_authors import load_document_authors

def main():
    print("Starting data loading into existing database...")
    
    try:
        
        if not os.path.exists(DB_PATH):
            print(f"Error: Database not found at {DB_PATH}")
            return
        
        print(f"Loading data into: {DB_PATH}")
        
        load_documents()
        load_authors() 
        load_organizations()
        load_document_authors()
        
        print("All data has been uploaded successfully!")
        
    except Exception as e:
        print(f"Fatal error: {e}")

if __name__ == "__main__":
    main()