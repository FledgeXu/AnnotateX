start-dev: 
	docker compose -p annotate-x-dev -f docker-compose-dev.yaml up -d --build --always-recreate-deps

stop-dev: 
	docker compose -p annotate-x-dev -f docker-compose-dev.yaml down

clean-dev:
	docker volume ls --filter name=annotate-x-dev -q | xargs -r docker volume rm

.PHONY: start-dev stop-dev clean-dev
