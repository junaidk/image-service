export POSTGRESQL_URL='postgres://app-user:secret@localhost:5432/app_db?sslmode=disable'

migrate -database ${POSTGRESQL_URL} -path migrations up