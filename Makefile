run: main
	./main

build:
	go build -o bin/app cmd/app/main.go

brun:
	go run cmd/app/main.go
