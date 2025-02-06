include .env
export

run_postgres:
	docker run --name auth-db-postgres -v /usr/share/auth-db -e POSTGRES_PASSWORD=zhabizhabi -p 5431:5432 postgres
.PHONY: run_postgres
compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down
