
docker-run:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=Pass1234 \
	  -e POSTGRES_USER=etzba  -e POSTGRES_DB=etzba -d postgres:14

sql-seed:
	PGPASSWORD=Pass1234 psql -h localhost -U etzba -d etzba < seed.sql

go-build:
	cd ../../cli && go build -o etz
	mv ../../cli/etz .

run-command:
	./etz sql --workers=3 --auth=secret.json --config=sql.csv
	./etz sql --workers=3 --auth=secret.json --config=sql.csv --duration=1s
	./etz sql --workers=3 --auth=secret.yaml --config=sql.csv --duration=1s --verbose=true
	./etz sql --workers=3 --auth=secret.yaml --config=sql.csv
