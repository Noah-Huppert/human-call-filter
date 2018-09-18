.PHONY: db migrate db-cli cloudsql-proxy

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
	psql -h ${DB_HOST} -U ${DB_USERNAME} -d ${DB_NAME}

# cloudsql-proxy starts a proxy which allows you to connect to the GCP Cloud 
# SQL instance
cloudsql-proxy:
	docker run \
		-it \
		--rm \
		--net host \
		-v ${PWD}/deploy/human-call-filter/postgres-client.service-account.json:/secrets/cloudsql/credentials.json \
		gcr.io/cloudsql-docker/gce-proxy:1.11 /cloud_sql_proxy \
			-instances=nh-k8s:us-central1:postgres=tcp:5432 \
			-credential_file=/secrets/cloudsql/credentials.json

