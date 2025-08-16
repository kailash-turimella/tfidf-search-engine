
# **TF-IDF Search Engine**

  

A high-performance search engine built in  **Go**  with  **SQLite**  backend support. It crawls, indexes, and ranks documents using the  **TF-IDF algorithm**  with support for stopword removal, stemming, and efficient query handling.

  

## **Features**

-   **Concurrent Crawling & Indexing**  – Built with goroutines, channels, and  sync.WaitGroup  for scalable document ingestion.
    
-   **Pluggable Index Interface** – Supports both **in-memory** and **SQLite** backends.
    
-   **Relational Schema Design**  – SQLite schema with normalized tables (Words,  URLs,  WordCounts) and foreign keys for fast lookups and referential integrity.
    
-   **TF-IDF Ranking**  – Computes term frequency–inverse document frequency scores to deliver relevant search results.
    
-   **Natural Language Processing**  – Includes stopword removal and Snowball stemming for higher retrieval precision.
    
-   **Comprehensive Testing**  – Unit and integration tests covering parsing, indexing, TF-IDF scoring, and concurrency.
    

  

## **Tech Stack**

-   **Language**: Go (Golang)
    
-   **Database**: SQLite
    
-   **Libraries**:
    -   go-sqlite3  – SQLite driver
    -   snowball  – stemming
    -   testify  – testing utilities
        
    

  

## **Project Structure**

```
├── html_files/
│  ├── index.html
│  ├── robots.txt
│  ├── crawl/
│  │  └── index.html
│  ├── search/
│  │  └── index.html
│  └── top10/
│  └── ... test files ...
├── clean.go
├── cleanHref_test.go
├── crawl.go
├── crawlDelay_test.go
├── crawl_test.go
├── database.go
├── disallow_test.go
├── download.go
├── download_test.go
├── extract.go
├── extract_test.go
├── hit.go
├── index.go  # Index  interface  definitions
├── indexDb.go  # SQLite/GORM-backed Index implementation
├── indexInmem.go  # In-memory Index implementation
├── main.go  # Entry point; starts HTTP server
├── robots.go  # robots.txt parsing & allow/deny + crawl delay
├── search_test.go
├── server.go  # HTTP handlers &  static  file serving
├── stop_test.go
├── stop.go
├── tfIdf_test.go
├── tfIdf.go
├── go.mod
├── go.sum
└── README.md
```



## **Getting Started**

  

### **1. Clone the Repository**

```
git clone git@github.com:kailash-turimella/tfidf-search-engine.git
cd tfidf-search-engine
```

### **2. Install Dependencies**
  

Ensure you have  [Go](https://go.dev/)  installed, then run:

```
go mod tidy
```

### **3. Run Tests**


Execute all unit and integration tests to verify functionality:

```
go test ./...
```

----------

## **Usage**

  

### **Start the Server**

```
go run .
```

### **Access the Web Interface**
  

Open your browser and navigate to the following endpoints:

-   **Home:**  [http://localhost:8080](http://localhost:8080/)
    
    
-   **Crawl UI:**  [http://localhost:8080/crawl](http://localhost:8080/crawl)
    
    -   Kick off a crawl manually, or start immediately with a URL:

```
http://localhost:8080/crawl?q=<start_url>
```


    
-   **Search UI:**  [http://localhost:8080/search](http://localhost:8080/search)
    
    -   Perform queries manually, or get ranked results directly:
        
    
```
http://localhost:8080/search?q=<term>
```
  

### **Index Modes**

-   -index=db  _(default)_  – Uses SQLite via GORM for persistent storage (database.db)
    
-   -index=inmem  – Uses an in-memory index (no persistence)
    

## **Schema Overview (SQLite)**

-   **Words**  – unique word tokens
    
-   **URLs**  – crawled document references
    
-   **WordCounts**  – occurrences of words in URLs with frequency counts
    

```sql
CREATE TABLE Words (
    id INTEGER PRIMARY KEY,
    word TEXT UNIQUE
);

CREATE TABLE URLs (
    id INTEGER PRIMARY KEY,
    url TEXT UNIQUE
);

CREATE TABLE WordCounts (
    word_id INTEGER,
    url_id INTEGER,
    count INTEGER,
    FOREIGN KEY (word_id) REFERENCES Words(id),
    FOREIGN KEY (url_id) REFERENCES URLs(id)
);
```
----------
