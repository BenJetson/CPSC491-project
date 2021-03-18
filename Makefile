
.PHONY: start
start:
	docker-compose up

.PHONY: clean
clean:
	docker-compose down -v
	yes | docker-compose rm

.PHONY: stop
stop:
	docker-compose stop

.PHONY: build
build: clean
	docker-compose build

.PHONY: reset-db
reset-db: stop
	docker-compose rm -s -f -v db
	docker volume rm cpsc491-project_persistent-dbstore
	@printf "\n🧼 DB has been erased. Re-run flyway to construct new DB.\n"

.PHONY: storybook
storybook:
	docker-compose up storybook
