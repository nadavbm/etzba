# this is a simple makefile to test etzba code changes 
# use it as a minimal ci before contributing 
#
# TAG docker image tag
TAG ?= latest
REPO ?= nadavbm/etzba
# all will run unit tests, build cli tool and run it with pgsql and api server
all: go-test go-build pgsql-up sql-seed run-pgsql-test pgsql-down api-up api-seed run-api-test api-down
pg: go-build pgsql-up sql-seed run-pgsql-test pgsql-down
api: go-build api-up api-seed run-api-test api-down
# go
.PHONY: go-test
go-test:
	go test -v ./...

.PHONY: go-build
go-build:
	cd cli && go build -o etz
	mv cli/etz .

# test postgers
pgsql-up:
	cd examples/pgsql && docker-compose down
	sleep 3
	cd examples/pgsql && docker-compose up -d

sql-seed:
	sleep 8
	PGPASSWORD=Pass1234 psql -h localhost -U etzba -d etzba < examples/pgsql/seed.sql

run-pgsql-test:
	ETZ_POSTGRES_DB=etzba ETZ_POSTGRES_HOST=localhost ETZ_POSTGRES_PORT=5432 ETZ_POSTGRES_USER=etzba ETZ_POSTGRES_PASSWORD=Pass1234 ./etz sql --workers=3 --config=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=3 --auth=examples/pgsql/secret.json --config=examples/pgsql/sql.csv --verbose
	./etz sql --workers=3 --auth=examples/pgsql/secret.json --config=examples/pgsql/sql.csv --verbose --rps=1
	./etz sql --workers=3 --auth=examples/pgsql/secret.json --config=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=10 --auth=examples/pgsql/secret.yaml --config=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=100 --auth=examples/pgsql/secret.yaml --config=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=20 --auth=examples/pgsql/secret.yaml --config=examples/pgsql/sql.csv --duration=10s --rps=50 --output=files/$$(date +%Y%m%d_%H%M%S)_result.json
	./etz sql --workers=200 --auth=examples/pgsql/secret.yaml --config=examples/pgsql/sql.csv --duration=10s --rps=500 --output=files/$$(date +%Y%m%d_%H%M%S)_result.json
	./etz sql --workers=400 --auth=examples/pgsql/secret.yaml --config=examples/pgsql/sql.csv --duration=10s --rps=1000 --output=files/$$(date +%Y%m%d_%H%M%S)_result.json

pgsql-down:
	cd examples/pgsql && docker-compose down
	sleep 5

# test api server
api-up:
	cd examples/api && docker-compose down
	sleep 3
	cd examples/api && docker-compose up -d	pg-database
	sleep 12
	cd examples/api && docker-compose up -d	etzba 

api-seed:
	cd examples/api && sh seed.sh

run-api-test:
	./etz api --workers=3 --auth=examples/api/secret.json --config=examples/api/api.json --verbose
	./etz api --workers=3 --auth=examples/api/secret.json --config=examples/api/api.json --verbose --rps=1
	./etz api --workers=3 --auth=examples/api/secret.json --config=examples/api/api.yaml --duration=3s
	./etz api --workers=10 --auth=examples/api/secret.yaml --config=examples/api/api.json --duration=3s
	./etz api --workers=100 --auth=examples/api/secret.yaml --config=examples/api/api.yaml --duration=3s
	./etz api --workers=10 --auth=examples/api/secret.yaml --config=examples/api/api.json --duration=10s --rps=50 --output=files/$$(date +%Y%m%d_%H%M%S)_result.json
	./etz api --workers=20 --auth=examples/api/secret.yaml --config=examples/api/api.yaml --duration=10s --rps=100 --output=files/$$(date +%Y%m%d_%H%M%S)_result.json
api-down:
	cd examples/api && docker-compose down

# prepare docker for cli tests
cleanup-docker: cleanup-api cleanup-pg

cleanup-api:
	docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=nadavbm/etzba-api-test:v0.0.1 --format="{{.ID}}"))

cleanup-pg:
	docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=postgres:14 --format="{{.ID}}"))

# build image and push to dockerhub
.PHONY: docker-build
docker-build:
	docker build -t ${REPO}:${TAG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${REPO}:${TAG}