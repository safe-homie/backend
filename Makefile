.PHONY: db dev clean

db:	
	@docker run --name postgres -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:17.2-alpine3.21

clean:
	@docker stop postgres && \
	docker rm postgres -v

dev: db
	@echo "[Dev ready]"

deploy:
	@docker compose --env-file .env.production up -d --build