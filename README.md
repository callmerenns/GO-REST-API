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
git clone https://github.com/callmerenns/GO-REST-API.git
cd GO-REST-API
```

### 2. Environment Configuration
Create a `.env` file in the root directory and configure your database and application settings:
```env
# Configuration Middleware
TOKEN_ISSUE=your_token_issue
TOKEN_SECRET=your_token_secret
TOKEN_EXPIRE=your_token_expire

# Configuration DB 
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_NAME=your_db_name
DB_DRIVER=your_db_driver

# Configuration APP
API_PORT=your_api_port
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
   docker build -t GO-REST-API .
   ```
2. Run the Docker container:
   ```bash
   docker run -p 8080:8080 GO-REST-API
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
| POST   | `/api/v1/auth/register`  | Register user         |
| POST   | `/api/v1/auth/login`     | Login user         |
| GET    | `/api/v1/auth/logout`    | Logout         |
| GET    | `/api/v1/products`       | Get all products         |
| GET    | `/api/v1/products/:id`   | Get a single product by id |
| GET    | `/api/v1/products/:id`   | Get a single product by stock |
| POST   | `/api/v1/products`       | Create a new product     |
| PUT    | `/api/v1/products/:id`   | Update an existing product |
| DELETE | `/api/v1/products/:id`   | Delete a product         |
| GET    | `/api/v1/profiles`       | Get all profiles         |
| GET    | `/api/v1/profiles/:id`   | Get a single profile by id |

### Example Request: Create Product
**POST** `/api/v1/products`
```json
{
    "status": {
        "code": 200,
        "message": "Ok"
    },
    "data": [
        {
            "ID": 1,
            "CreatedAt": "2024-12-25T21:48:31.698+07:00",
            "UpdatedAt": "2024-12-25T21:48:31.698+07:00",
            "DeletedAt": {
                "Time": "0001-01-01T00:00:00Z",
                "Valid": false
            },
            "name": "Shoes",
            "description": "Shoes H&M for Women's Fashion",
            "stock": 16,
            "price": 1323316,
            "users": [
                {
                    "ID": 1,
                    "CreatedAt": "2024-12-25T21:48:05.611+07:00",
                    "UpdatedAt": "2024-12-25T21:48:05.601+07:00",
                    "DeletedAt": {
                        "Time": "0001-01-01T00:00:00Z",
                        "Valid": false
                    },
                    "firstname": "Paul",
                    "lastname": "Casey",
                    "email": "paul.casey.1@gslingacademy.com",
                    "role": "admin"
                }
            ]
        }
    ],
    "paging": {
        "page": 1,
        "rowsPerPage": 10,
        "totalRows": 1,
        "totalPages": 1
    }
}
```

---

## Testing
Use tools like [Postman](https://www.postman.com/) or [cURL](https://curl.se/) to test the endpoints.

### Example cURL Commands
- **Get All Products**:
  ```bash
  curl -X GET http://localhost:8080/api/v1/products
  ```
- **Create Product**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"name":"Sample Product","description":"Sample Description","stock":10,"price":100.0}' http://localhost:8080/api/v1/products
  ```

---

## Project Structure
```
.
├── cmd
│   ├── config        # Configuration files
│   ├── delivery      # HTTP handlers 
|       ├── controllers        # Configuration enpoint url
|       ├── middlewares        # Configuration middleware
|       └── server.go          # Entry point of the application
│   ├── entity        # Domain entities and models
|       ├── dto       # Data transfer object
│   ├── repository    # Data access logic
|   ├── shared        # Shared utilities and helpers
|       ├── common             # Custom response
|       ├── model              # Model for response data
|       └── service            # Configuration JWT
│   ├── usecase       # Business logic
│   └── utils         # Miscellaneous utilities
├── docs              # Configuration swagger
├── .air.toml
├── .env              # Environment file
├── .env.example      # Example environment file
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── main.go           # Entry point of the application
```

---

## License
This project is licensed under the [MIT License](LICENSE).

---

## Contact
For questions or suggestions, please contact:
- **Email**: altsaqifnugraha19@gmail.com
- **GitHub**: [callmerenns](https://github.com/callmerenns)

---

Happy Coding! 🚀
