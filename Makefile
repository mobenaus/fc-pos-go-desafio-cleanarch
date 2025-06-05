BINARY_NAME=ordersystem
SHELL := bash
 

build:
	cd cmd/ordersystem && go build -o ${BINARY_NAME} main.go wire_gen.go

run: services
	cd cmd/ordersystem && go run main.go wire_gen.go
 
test:
	go clean -testcache
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