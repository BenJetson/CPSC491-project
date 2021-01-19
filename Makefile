
.PHONY: start
start:
	docker-compose up

.PHONY: clean
clean:
	docker-compose down -v
	yes | docker-compose rm

.PHONY: stop
stop:
	docker-compose down

.PHONY: build
build: clean
	docker-compose build

