# REST API CRUD Product

A simple REST API for managing products, built using:
- **Framework**: Gin (Golang)
- **Database**: MySQL
- **ORM**: GORM
- **Concurrency**: Goroutines
- **Containerization**: Docker
- **Architecture**: Clean Code Architecture

## Features
- Create, Read, Update, and Delete (CRUD) operations for products.
- Clean and maintainable codebase following Clean Code principles.
- Supports concurrency for efficient data processing.

---

## Requirements
Ensure you have the following installed:
- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [MySQL](https://www.mysql.com/)

---

## Installation and Setup
Follow these steps to run the project:

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/crud-product-api.git
cd crud-product-api
```

### 2. Environment Configuration
Create a `.env` file in the root directory and configure your database and application settings:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=productdb
APP_PORT=8080
```

### 3. Build and Run Using Docker
#### Using Docker Compose:
1. Ensure `docker-compose.yml` is properly configured.
2. Start the containers:
   ```bash
   docker-compose up --build
   ```
3. The application will be available at `http://localhost:8080`.

#### Without Docker Compose:
1. Build the Docker image:
   ```bash
   docker build -t crud-product-api .
   ```
2. Run the Docker container:
   ```bash
   docker run -p 8080:8080 crud-product-api
   ```

### 4. Run Locally Without Docker
1. Install dependencies:
   ```bash
   go mod tidy
   ```
2. Start the server:
   ```bash
   go run main.go
   ```
3. The application will run at `http://localhost:8080`.

---

## API Endpoints

| Method | Endpoint          | Description              |
|--------|-------------------|--------------------------|
| GET    | `/products`       | Get all products         |
| GET    | `/products/:id`   | Get a single product     |
| POST   | `/products`       | Create a new product     |
| PUT    | `/products/:id`   | Update an existing product |
| DELETE | `/products/:id`   | Delete a product         |

### Example Request: Create Product
**POST** `/products`
```json
{
  "name": "Sample Product",
  "price": 100.0,
  "quantity": 10
}
```

---

## Testing
Use tools like [Postman](https://www.postman.com/) or [cURL](https://curl.se/) to test the endpoints.

### Example cURL Commands
- **Get All Products**:
  ```bash
  curl -X GET http://localhost:8080/products
  ```
- **Create Product**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"name":"Sample Product","price":100.0,"quantity":10}' http://localhost:8080/products
  ```

---

## Project Structure
```
.
├── cmd
│   ├── config        # Configuration files
│   ├── delivery      # HTTP handlers (controllers)
│   ├── entity        # Domain entities and models
│   ├── repository    # Data access logic
│   ├── usecase       # Business logic
│   ├── shared        # Shared utilities and helpers
│   └── utils         # Miscellaneous utilities
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go           # Entry point of the application
└── .env.example      # Example environment file
```

---

## License
This project is licensed under the [MIT License](LICENSE).

---

## Contact
For questions or suggestions, please contact:
- **Email**: your.email@example.com
- **GitHub**: [your-username](https://github.com/your-username)

---

Happy Coding! 🚀
