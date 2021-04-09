
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
	docker volume rm cpsc491_persistent-dbstore
	@printf "\n🧼 DB has been erased. Re-run flyway to construct new DB.\n"

.PHONY: stop-clean-testdb
stop-clean-testdb:
	docker-compose stop flyway-testdb testdb
	docker-compose rm -f flyway-testdb testdb

.PHONY: testdb
testdb: stop-clean-testdb
	docker-compose up flyway-testdb testdb

.PHONY: testdb-background
testdb-background: stop-clean-testdb
	docker-compose up -d flyway-testdb testdb

.PHONY: testgo
testgo:
	@./go/scripts/test.sh

.PHONY: storybook
storybook:
	docker-compose up storybook
