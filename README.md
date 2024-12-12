# Book Service

The **Book Service** is a highly modular and scalable backend application built in Go, designed to manage books and their metadata efficiently. It leverages modern web frameworks, Elasticsearch, Redis, and follows best practices for RESTful APIs.


## 🚀 Features

- **Book Management API:**
    - Add, update, delete, search, and retrieve book details.
- **Middleware Support:**
    - Logging, request validation, and user action tracking.
- **Elasticsearch Integration:**
    - Indexing and querying books with custom mappings.
- **Redis Integration:**
    - Buffering and caching for improved performance.
- **Modular Design:**
    - Organized packages with clear separation of concerns.

---

## 🛠️ Setup and Installation

### Prerequisites

- **Go** (>= 1.17)
- **Docker** (for running Elasticsearch and Redis locally)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/moshe-mizrachi/book_service.git
   cd book_service

2. Install Dep
    ```bash
    go mod tidy
    ```
3. Run the app
    ```bash
    go run cmd/main/main.go
    ```

### 🔑 Environment Variables
Define the following variables in a .env or .env.test file:

| Variable  | Description                      | Default               |
|-----------|----------------------------------|-----------------------|
| PORT      | Port for the application server  | 8080                 |
| ELS_URI   | Elasticsearch connection URI     | http://localhost:9200 |
| REDIS_URI | Redis connection URI             | localhost:6379        |

### 📖 API Endpoints
Books API (/v1/books)

| Method    | Endpoint       | Description                  |
|-----------|----------------|------------------------------|
| `POST`    | `/v1/books/`   | Add a new book              |
| `GET`     | `/v1/books/:id`| Retrieve book details by ID |
| `PUT`     | `/v1/books/:id`| Update book title by ID     |
| `DELETE`  | `/v1/books/:id`| Delete a book by ID         |
| `GET`     | `/v1/books/search` | Search for books          |


### Middlewares

```go
// Logger Middleware
Logs details of incoming requests, including:
- HTTP method
- Request path
- Status code
- Latency
- User agent

// Validation Middleware
Validates incoming requests against defined schemas using Gin's binding and validation mechanisms.

// RecordActions Middleware
Tracks user actions such as method, path, and time, and stores them in Redis for auditing purposes.
```


### Query Builder
The query builder lets you play and construct the wanted query according to your needs

```go
// Query Builder Methods
qb := query.NewQueryBuilder()

// Add an ID filter to the query
qb.ID("12345")

// Add a title filter to the query
qb.Title("Book Title")

// Add an author name filter to the query
qb.AuthorName("Author Name")

// Add a price range filter to the query
qb.PriceRange(10, 50)

// Build the final query
result := qb.Build()

// Output Example:
{
  "query": {
    "bool": {
      "must": [
        {"term": {"_id": "12345"}},
        {"match": {"title": "Book Title"}}
      ],
      "filter": [
        {"range": {"price": {"gte": 10, "lte": 50}}}
      ]
    }
  }
}
```

### 🧪 Testing
#### Run Tests
To run the tests:

```bash
  go test ./test/... -v
```
#### Tests include:
1. Unit tests for query builders.
2. Integration tests for API endpoints.

## 📦 Built With

- **Gin** - Web framework
- **Elasticsearch** - Full-text search engine
- **Redis** - In-memory data store
- **Logrus** - Logging library


## 👨‍💻 Contributors

- **Moshe Mizrachi Fiverr** - Initial work

Feel free to contribute by submitting issues or pull requests.
