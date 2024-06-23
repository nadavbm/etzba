docker-run:
	docker-compose down
	docker-compose up -d

go-build:
	cd ../../cli && go build -o etz
	mv ../../cli/etz .

run-command:
	./etz api --workers=3 --auth=secret.json --config=api.json
	./etz api --workers=3 --auth=secret.json --config=api.json --duration=1s
	./etz api --workers=3 --auth=secret.yaml --config=api.yaml --duration=1s --verbose=true
	./etz api --workers=3 --auth=secret.yaml --config=api.yaml
