# Golang todo list (TBD)

This is the repository for my first golang API project.

## project features
- Built in Go version 1.23
- Run through Docker & docker compose
  - docker version 28.0
- Uses the standard `net/http` library
- Uses `go-sql-driver/mysql` as my DB driver
- Uses `golang-jwt` for authentication
- Uses `joho/godotenv` for loading env file
- Uses `golang.org/x/crypto` for password hashing
- Uses `golang-migrate/migrate/v4` for DB migration
- Uses `go-playground/validator/v10` for request payload validation
- MORE TO COME...

## Run
- copy your own .env file: `cp .env.example .env`
- edit your env variables (enclosed within square brackets)
```sh
ENV=production  # Please remain this variable same to not use godotenv in docker container 
JWT_SECRET=[your_jwtsecret]
MYSQL_DSN=root:[your_password]@(my-mysql:3306)/[your_dbname]?parseTime=true
MYSQL_ROOT_PASSWORD=[your_password]
MYSQL_DATABASE=[your_dbname]
```
- Run docker compose: `docker compose up -d`
- The port is listening on port 11451. That'll do it.