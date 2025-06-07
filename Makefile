BINARY_NAME=ordersystem
SHELL := bash
 

build:
	cd cmd/ordersystem && go build -o ${BINARY_NAME} main.go wire_gen.go

run: services
	go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go 
 
test:
	go test ./...
 
clean:
	go clean
	rm ${BINARY_NAME}

services:
	docker compose up -d
	@echo "Wait for rabbitmq started..."
	@d=1
	@while [ "$${d}" != "0" ]; do \
		docker compose logs rabbitmq 2>&1 | grep -o "Server startup complete"; \
		d="$$?" ; \
		sleep 3; \
	done

services-shutdown:
	docker compose down

services-wipe: services-shutdown
	sudo rm -Rf .docker

generate-graphql:
	go run github.com/99designs/gqlgen generate

generate-grpc:
	protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

wire-di:
	wire
