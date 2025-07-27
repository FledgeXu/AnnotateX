start-dev: 
	docker compose -p annotate-x-dev -f docker-compose-dev.yaml up -d --build --always-recreate-deps

stop-dev: 
	docker compose -p annotate-x-dev -f docker-compose-dev.yaml down

clean-dev:
	docker volume ls --filter name=annotate-x-dev -q | xargs -r docker volume rm

start-backend: 
	docker compose -p annotate-x-backend -f docker-compose-backend.yaml up -d --build --always-recreate-deps

stop-backend: 
	docker compose -p annotate-x-backend -f docker-compose-backend.yaml down

clean-backend:
	docker volume ls --filter name=annotate-x-backend -q | xargs -r docker volume rm

.PHONY: start-dev stop-dev clean-dev start-backend stop-backend clean-backend
