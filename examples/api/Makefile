docker-run:
	docker-compose down
	docker-compose up -d

go-build:
	cd ../../cli && go build -o etz
	mv ../../cli/etz .

run-command:
	./etz api --workers=3 --config=secret.json --helpers=api.json
	./etz api --workers=3 --config=secret.json --helpers=api.json --duration=1s
