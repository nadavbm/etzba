docker-run:
	docker-compose down
	docker-compose up -d

go-build:
	cd ../../cli && go build -o etz
	mv ../../cli/etz .

run-command:
	./etz api --workers=3 --config=secret.json --helpers=api.json
	./etz api --workers=3 --config=secret.json --helpers=api.json --duration=1s
	./etz api --workers=3 --config=secret.yaml --helpers=api.yaml --duration=1s --verbose=true
	./etz api --workers=3 --config=secret.yaml --helpers=api.yaml
