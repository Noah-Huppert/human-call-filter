.PHONY: db migrate db-cli

NAME=dev-human-call-filter

# db starts a local database
db:
	docker run \
		-it \
		--rm \
		--net host \
		-v ${PWD}/run-data:/var/lib/postgresql/data \
		-e POSTGRES_DB=${NAME} \
		-e POSTGRES_USER=${NAME} \
		postgres

# migrate runs database migrations
migrate:
	go run scripts/db-migrate.go

# db-cli connects to the local database with the psql command line interface
db-cli:
	psql -h localhost -U ${USER} -d ${NAME}
