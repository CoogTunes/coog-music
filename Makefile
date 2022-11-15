run:
	go run ./cmd/web/.
migrateup:
	migrate -path db/migration -database "postgresql://postgres:coog2022@coogtune.ccpw7qggmv2b.us-east-2.rds.amazonaws.com:5433/music1?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:coog2022@coogtune.ccpw7qggmv2b.us-east-2.rds.amazonaws.com:5433/music1?sslmode=disable" -verbose down
.PHONY: run migrateup migratedown