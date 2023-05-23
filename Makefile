# this is a simple makefile to test etzba code changes 
# use it as a minimal ci before contributing 
#
# all will run unit tests, build cli tool and run it with pgsql and api server
all: go-test go-build pgsql-up sql-seed run-pgsql-test pgsql-down api-up api-seed run-api-test api-down 
# go
go-test:
	go test -v ./...

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
	./etz sql --workers=3 --config=examples/pgsql/secret.json --helpers=examples/pgsql/sql.csv --verbose
	./etz sql --workers=3 --config=examples/pgsql/secret.json --helpers=examples/pgsql/sql.csv --verbose --rps=1
	./etz sql --workers=3 --config=examples/pgsql/secret.json --helpers=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=10 --config=examples/pgsql/secret.yaml --helpers=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=100 --config=examples/pgsql/secret.yaml --helpers=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=10 --config=examples/pgsql/secret.yaml --helpers=examples/pgsql/sql.csv --duration=10s --rps=5
	./etz sql --workers=100 --config=examples/pgsql/secret.yaml --helpers=examples/pgsql/sql.csv --duration=10s --rps=10

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
	./etz api --workers=3 --config=examples/api/secret.json --helpers=examples/api/api.json --verbose
	./etz api --workers=3 --config=examples/api/secret.json --helpers=examples/api/api.json --verbose --rps=1
	./etz api --workers=3 --config=examples/api/secret.json --helpers=examples/api/api.yaml --duration=3s
	./etz api --workers=10 --config=examples/api/secret.yaml --helpers=examples/api/api.json --duration=3s
	./etz api --workers=100 --config=examples/api/secret.yaml --helpers=examples/api/api.yaml --duration=3s
	./etz api --workers=10 --config=examples/api/secret.yaml --helpers=examples/api/api.json --duration=10s --rps=200
	./etz api --workers=100 --config=examples/api/secret.yaml --helpers=examples/api/api.yaml --duration=10s -- rps=400

api-down:
	cd examples/api && docker-compose down

# prepare docker for cli tests
cleanup-docker: cleanup-api cleanup-pg

cleanup-api:
	docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=nadavbm/etzba-api-test:v0.0.1 --format="{{.ID}}"))

cleanup-pg:
	docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=postgres:14 --format="{{.ID}}"))