.PHONY: db migrate

# db starts a local database
db:
	docker run \
		-it \
		--rm \
		--net host \
		-v ${PWD}/run-data:/var/lib/postgresql/data \
		-e POSTGRES_DB=dev-human-call-filter \
		-e POSTGRES_USER=dev-human-call-filter \
		postgres

# migrate runs database migrations
migrate:
	go run scripts/db-migrate.go
