.PHONY: dev up down restart logs migrate-up migrate-down psql clean db-only

# Start development dengan docker compose
dev:
	docker-compose up -d

# Naikan semua container
up:
	docker-compose up -d

# Hanya menjalankan database dan pgadmin
db-only:
	docker-compose up -d postgres pgadmin

# Matikan semua container
down:
	docker-compose down

# Restart semua container
restart:
	docker-compose down
	docker-compose up -d

# Melihat logs
logs:
	docker-compose logs -f

# Menjalankan migrasi database (UP)
migrate-up:
	docker-compose exec api migrate -path=./migrations -database "postgres://postgres:postgres@postgres:5432/launcher_db?sslmode=disable" up

# Rollback migrasi database (DOWN)
migrate-down:
	docker-compose exec api migrate -path=./migrations -database "postgres://postgres:postgres@postgres:5432/launcher_db?sslmode=disable" down 1

# Masuk ke postgres CLI
psql:
	docker-compose exec postgres psql -U postgres -d launcher_db

# Membersihkan volume dan image yang tidak terpakai
clean:
	docker-compose down -v
	docker system prune -f 