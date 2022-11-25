run:
	go run ./cmd/web/.
migrateup:
	migrate -path db/migration -database "postgresql://postgres:ivan3792@localhost:5433/music?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:ivan3792@localhost:5433/music?sslmode=disable" -verbose down
.PHONY: run migrateup migratedown