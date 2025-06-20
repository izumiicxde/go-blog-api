# Go-Blog

A scalable RESTful Blogging API built using Go (Golang), leveraging Gorilla Mux for routing, GORM as the ORM, and PostgreSQL as the database.

## Features

- User authentication with JWT
- Email validation and verification with OTP
- Create, Read, Update, Delete (CRUD) blog posts
- PostgreSQL as the persistent storage
- Input validation and error handling
- Modular folder structure

## API Endpoints

### Authentication

- POST /register - Register a new user
- POST /login - User login
- POST /verify - Verify user (using code/otp)
- GET /get-verification-code - Send verification code to the user's email

### Blog Operations [Must be authenticated]

"PUBLIC"

- GET /blogs - Fetch all blog posts

- GET /blogs/{userId} - Get blogs by userId
- POST /blogs - Create a new blog post
- GET /blogs/{id} - Fetch a single blog post by ID
- PATCH /blogs/{id} - Update a blog post by ID
- DELETE /blogs/soft/{id} - Soft delete a blog post (mark as deleted)
- DELETE /blogs/delete/{id} - Hard delete a blog post (remove permanently)

## Built With

- [Gorilla Mux](https://github.com/gorilla/mux) — HTTP request router and dispatcher for building Go web servers.
- [GORM](https://gorm.io/) — ORM library for Go, used for interacting with the PostgreSQL database.
- [PostgreSQL](https://www.postgresql.org/) — Relational database used for storing user and blog data.
- [godotenv](https://github.com/joho/godotenv) — Loads environment variables from `.env` files.
- [JWT (github.com/golang-jwt/jwt/v5)](https://github.com/golang-jwt/jwt) — Used for implementing JSON Web Token-based authentication.
- [Gomail](https://github.com/go-gomail/gomail) — Package used to send emails (for verification codes).
- [Golang-Migrate](https://github.com/golang-migrate/migrate) — Database migration tool used to handle schema migrations in PostgreSQL.

## Run Locally

First clone this repo with

```bash
$ git clone https://github.com/izumiicxde/go-blog-api.git
```

Then set the environmental variables

```env
PUBLIC_HOST="http://localhost:8080"
PORT="8080"

DB_NAME=""
DB_HOST=""
DB_USER=""
DB_PASSWORD=""
DB_PORT="5432"
DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

JWT_SECRET=""
JWT_EXPIRATION=

# Gomail configuration
SMTP_SERVER=smtp.example.com
SMTP_PORT=
SMTP_USER=something@example.com
SMTP_PASSWORD=""


```

```bash
$ go mod tidy
$ make run
```

This will start the local server at `http://localhost:8080/api/v1`

---

## Acknowledgements

- [Gorilla Mux](https://github.com/gorilla/mux)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)

## License

[MIT](https://choosealicense.com/licenses/mit/)
