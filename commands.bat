@echo off

if "%1"=="dev" (
    docker-compose up -d
    goto :eof
)

if "%1"=="up" (
    docker-compose up -d
    goto :eof
)

if "%1"=="db-only" (
    docker-compose up -d postgres pgadmin
    goto :eof
)

if "%1"=="down" (
    docker-compose down
    goto :eof
)

if "%1"=="restart" (
    docker-compose down
    docker-compose up -d
    goto :eof
)

if "%1"=="logs" (
    docker-compose logs -f
    goto :eof
)

if "%1"=="migrate-up" (
    docker-compose exec api migrate -path=./migrations -database "postgres://postgres:postgres@postgres:5432/launcher_db?sslmode=disable" up
    goto :eof
)

if "%1"=="migrate-down" (
    docker-compose exec api migrate -path=./migrations -database "postgres://postgres:postgres@postgres:5432/launcher_db?sslmode=disable" down 1
    goto :eof
)

if "%1"=="psql" (
    docker-compose exec postgres psql -U postgres -d launcher_db
    goto :eof
)

if "%1"=="clean" (
    docker-compose down -v
    docker system prune -f
    goto :eof
)

echo Command tidak ditemukan 