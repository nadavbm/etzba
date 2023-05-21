# this is a simple makefile to test etzba code changes 
# use it as a minimal ci before contributing 
#
# all will run unit tests, build cli tool and run it with pgsql and api server
all: go-test go-build pgsql-up sql-seed run-pgsql-test pgsql-down api-up api-seed run-api-test api-downdocker 

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
	./etz sql --workers=3 --config=examples/pgsql/secret.json --helpers=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=10 --config=examples/pgsql/secret.json --helpers=examples/pgsql/sql.csv --duration=3s
	./etz sql --workers=100 --config=examples/pgsql/secret.json --helpers=examples/pgsql/sql.csv --duration=3s

pgsql-down:
	cd examples/pgsql && docker-compose down
	sleep 5

# test api server
api-up:
	cd examples/api && docker-compose down
	sleep 3
	cd examples/api && docker-compose up -d	

api-seed:
	sleep 12
	curl -X POST http://localhost:8080/v1/signup \
	   -H 'Content-Type: application/json' \
	   -d '{"name": "etzba","email": "etzba@etzba.com","password": "Pass1234"}' | jq '.token'

run-api-test:
	./etz api --workers=3 --config=examples/api/secret.json --helpers=examples/api/api.json --verbose
	./etz api --workers=3 --config=examples/api/secret.json --helpers=examples/api/api.json --duration=3s
	./etz api --workers=10 --config=examples/api/secret.json --helpers=examples/api/api.json --duration=3s
	./etz api --workers=100 --config=examples/api/secret.json --helpers=examples/api/api.json --duration=3s

api-down:
	cd examples/api && docker-compose down