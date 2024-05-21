start:
	go build -o ./app ./cmd/app/main.go
	./app
	rm ./app


migrate:
	go build -o ./mig ./cmd/migrations/main.go
	./mig
	rm ./mig
