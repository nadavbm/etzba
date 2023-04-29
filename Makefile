go-build:
	cd cli && go build -o etz

pgsql-run:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=Pass1234 \
	  -e POSTGRES_USER=etzba  -e POSTGRES_DB=etzba -d postgres:14

sql-seed:
	PGPASSWORD=Pass1234 psql -h localhost -U etzba -d etzba < examples/pgsql/seed.sql