postgres:
	docker compose up -d --build postgres

stop-postgres:
	docker stop image-processing-postgres

createdb:
	docker exec -it image-processing-postgres createdb --username=postgres --owner=postgres image_processing

dropdb:
	docker exec -it image-processing-postgres dropdb --username=postgres image_processing

.PHONY: postgres stop-postgres createdb dropdb