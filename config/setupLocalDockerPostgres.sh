docker run -d \
    --name sample-postgres \
    -e POSTGRES_PASSWORD=sample \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -e POSTGRES_DB=sample \
    -e POSTGRES_USER=sample \
    -p 5432:5432 \
    postgres:12.3