include .env
export

compose-up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up
compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down
