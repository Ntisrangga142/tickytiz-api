include ./.env
DBURL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONPATH=db/migrations
SEEDSPATH=db/seed

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONPATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) up

migrate-down:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) down $(s)

migrate-status:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) version

migrate-force:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) force $(v)

seed:
	for file in $(SEEDSPATH)/*.sql; do \
		psql $(DBURL) -f $$file; \
	done
