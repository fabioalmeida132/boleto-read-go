.PHONY: start-dev
start-dev:
	docker-compose up -d

.PHONY: stop-dev
stop-dev:
	docker-compose down

.PHONY: restart-dev
restart-dev:
	docker-compose down && docker-compose up -d