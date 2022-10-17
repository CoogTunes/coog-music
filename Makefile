run:
	go run ./cmd/web/.
migrateup:
	migrate -path db/migration -database "postgresql://root:dummypassword@localhost:5432/music?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:dummypassword@localhost:5432/music?sslmode=disable" -verbose down
.PHONY: run migrateup migratedown