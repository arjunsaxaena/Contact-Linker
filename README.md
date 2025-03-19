# Moonrider Assignment

This project is a contact management system that allows users to identify and manage contacts based on email and phone numbers. It uses a PostgreSQL database for storing contact information and provides a REST API for interaction.

## Features

- Create and manage contacts with email and phone numbers.
- Automatically link contacts with the same email or phone number.
- Distinguish between primary and secondary contacts.
- REST API for identifying and managing contacts.

## Technologies Used

- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Framework**: Gin (for REST API)
- **Environment Management**: `godotenv`
- **SQL Builder**: `go-sqlbuilder`
- **Containerization**: Docker

## Prerequisites

- Go 1.20 or later
- PostgreSQL
- Docker (optional, for containerized deployment)

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/arjunsaxaena/Moonrider-Assignment.git
cd Moonrider-Assignment
```

### 2. Configure Environment Variables

Create a `.env` file in the root directory with the following content:

```properties
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=moonrider
DB_SSLMODE=disable
```

### 3. Run Database Migrations

Ensure PostgreSQL is running and execute the SQL migration script:

```bash
psql -U postgres -d moonrider -f migrations/1_create_table.sql
```

### 4. Run the Application

#### Using Go

```bash
go run cmd/main.go
```

#### Using Docker

Build and run the Docker container:

```bash
docker build -t moonrider-assignment .
docker run -p 8080:8080 --env-file .env moonrider-assignment
```

### 5. Test the API

Use tools like Postman or `curl` to test the API. Example:

```bash
curl -X POST http://localhost:8080/identify \
-H "Content-Type: application/json" \
-d '{"email": "test@example.com", "phoneNumber": "1234567890"}'
```

## API Endpoints

### POST `/identify`

**Request Body**:
```json
{
  "email": "test@example.com",
  "phoneNumber": "1234567890"
}
```

**Response**:
```json
{
  "primaryContactId": "uuid",
  "emails": ["test@example.com"],
  "phoneNumbers": ["1234567890"],
  "secondaryContactIds": []
}
```

## Project Structure

```
Moonrider-Assignment/
├── cmd/                # Main entry point
├── config/             # Configuration loading
├── controllers/        # API controllers
├── model/              # Data models
├── repository/         # Database operations
├── migrations/         # Database migration scripts
├── Dockerfile          # Docker configuration
├── .env                # Environment variables
└── README.md           # Project documentation
```

## Author

Arjun Saxena