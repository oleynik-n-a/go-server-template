# go-server-template
Server template that you can use for your projects

## Tech stack
- Postgres
- Go 1.24.1
- Gin
- Zap logger
- Open API 3.0
- JWT auth
- Docker

## Setup
`.env` file should be located in source dir (`/go-server-template`, not `/src`).

Example of `.env` file:
```env
DB_HOST=db
DB_PORT=5432
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=yourdb
SSL_MODE=disable

PASSWORD_SALT=yourpasswordsalt
JWT_ACCESS_SECRET=yourjwtaccesssecret
JWT_REFRESH_SECRET=yourjwtrefreshsecret
```

## Usage
After adding .env file you can clone project, open source dir and run `docker-compose up --build` terminal.
