# STOQR API

```bash
docker run --name stoqr-postgres -p 5432:5432 -e POSTGRES_PASSWORD=stoqr -d postgres

export STOQR_API_DB_HOST=localhost
export STOQR_API_DB_USER=postgres
export STOQR_API_DB_PASSWORD=stoqr
export STOQR_API_DB_NAME=postgres
export STOQR_API_DB_PORT=5432

go build .
./stoqr-api
```