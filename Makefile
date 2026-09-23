postgres:
	docker compose up -d --build postgres

stop-postgres:
	docker stop image-processing-postgres

createdb:
	docker exec -it image-processing-postgres createdb --username=postgres --owner=postgres image_processing

dropdb:
	docker exec -it image-processing-postgres dropdb --username=postgres image_processing

migrateup:
	migrate -path database/migration -database "postgres://postgres:postgres@localhost:5433/image_processing?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migration -database "postgres://postgres:postgres@localhost:5433/image_processing?sslmode=disable" -verbose down

.PHONY: postgres stop-postgres createdb dropdb