.PHONY: db timescaledb dev clean

db:	
	@docker run --name postgres -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:17.2-alpine3.21

timescaledb:
	@docker run --name timescaledb -d -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb:2.18.1-pg17

clean:
	@docker stop timescaledb && \
	docker rm timescaledb -v

dev: timescaledb
	@echo "[Dev ready]"

deploy:
	@docker compose --env-file .env.production up -d --build