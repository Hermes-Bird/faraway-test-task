server_address = pow-server:8080

run-server:
	go run cmd/server/main.go
build-server:
	go build -o app ./cmd/server/main.go
build-client:
	go build -o app ./cmd/client/main.go

create-network:
	docker network create pow-server-network

build-docker-server:
	docker build -t pow-server -f DockerfileServer .
run-docker-server:
	docker run --network=pow-server-network --name=pow-server pow-server

build-docker-client:
	docker build -t pow-client -f DockerfileClient .
run-docker-client:
	docker run --network=pow-server-network pow-client ./app $(server_address)
