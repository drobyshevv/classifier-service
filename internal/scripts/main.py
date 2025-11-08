import sys
import os

# import config
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(__file__))))

from config import *
from load_documents import load_documents
from load_authors import load_authors
from load_organizations import load_organizations
from load_document_authors import load_document_authors

def main():
    print("the start of uploading data to SQLite...")
    
    try:
        load_documents()
        load_authors() 
        load_organizations()
        load_document_authors()  #foreign keys
        
        print("all data has been uploaded successfully!")
        
    except Exception as e:
        print(f"fatal error: {e}")

if __name__ == "__main__":
    main()