build-local:
	pwd
	go build -o app cmd/main.go

build-alpine:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/main.go

run:
	make build-alpine
	docker-compose up -d

teardown:
	docker-compose down
