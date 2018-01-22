NAME = ppp
DB_DOCKER_CONTAINER = db-dev
export PGPASSWORD=mysecretpassword

# This target ensures that the pre-commit hook is installed and kept up to date
# if pre-commit updates.
pre-commit: .git/hooks/pre-commit
.git/hooks/pre-commit: /usr/local/bin/pre-commit
	pre-commit install

deps: pre-commit client_deps server_deps
test: client_test server_test

client_deps_update:
	yarn upgrade
client_deps: .client_deps.stamp
.client_deps.stamp: yarn.lock
	yarn install
	touch .client_deps.stamp
client_build: client_deps
	yarn build
client_run: client_deps
	yarn start
client_test: client_deps
	yarn test

server_deps_update:
	dep ensure -update
server_deps: .server_deps.stamp
.server_deps.stamp: Gopkg.lock
	bin/check_gopath.sh
	dep ensure
	go install ./vendor/github.com/markbates/pop/soda
	go install ./vendor/github.com/golang/lint/golint
	touch .server_deps.stamp
server_build: server_deps
	go build -io bin/webserver ./cmd/webserver
server_run_only: db_dev_run
	./bin/webserver \
		-port :8080 \
		-debug_logging
server_run: client_build server_build server_run_only
server_run_dev: server_build server_run_only

server_build_docker:
	docker build . -t ppp:dev
server_run_only_docker: db_dev_run
	docker run --name ppp -p 8080:8080 ppp:dev

server_test: db_dev_run db_test_reset
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		go test ./...

db_dev_init:
	docker run --name $(DB_DOCKER_CONTAINER) \
		-e \
		POSTGRES_PASSWORD=$(PGPASSWORD) \
		-d \
		-p 5432:5432 \
		postgres:latest
	bin/wait-for-db
	createdb -p 5432 -h localhost -U postgres dev_db
db_dev_run:
	# We don't want to utilize Docker to start the database if we're
	# in the CircleCI environment. It has its own configuration to launch
	# a DB.
	[ -z "$(CIRCLECI)" ] && \
		docker start $(DB_DOCKER_CONTAINER) || \
		echo "Relying on CircleCI's database container."
db_dev_reset:
	echo "Attempting to reset local dev database..."
	docker kill $(DB_DOCKER_CONTAINER) &&	\
		docker rm $(DB_DOCKER_CONTAINER) || \
		echo "No dev database"
db_dev_migrate: db_dev_run
	soda migrate up
db_dev_migrate_down: db_dev_run
	soda migrate down

db_test_reset:
	# Initialize a test database if we're not in a CircleCI environment.
	[ -z "$(CIRCLECI)" ] && \
		dropdb -p 5432 -h localhost -U postgres --if-exists test_db && \
		createdb -p 5432 -h localhost -U postgres test_db || \
		echo "Relying on CircleCI's test database setup."
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		bin/wait-for-db
	soda -e test migrate up

.PHONY: pre-commit deps test client_deps client_build client_run client_test
.PHONY: server_deps_update server_deps server_build server_run_only server_run server_run_dev server_build_docker server_run_only_docker server_test
.PHONY: db_dev_init db_dev_run db_dev_reset db_dev_migrate db_dev_migrate_down db_test_reset
