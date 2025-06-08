# fc-pos-go-desafio-cleancode

## Descrição do Desafio
```
Desafio Clean Architecture

Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.
```
## TL;DR;

Requisitos:
- docker
- make
- go
- evans (para acessar via gRPC)
- browser (para acessar via GraphQL)
- VS-Code com plugin "REST Client" (para acessar via chamadas API REST)
  - arquivo:
    -  api/create_order.http
    -  api/list_orders.http

Execucao:
- make run


## Servicos e portas

Por default o servico responde os protocolos nas seguintes portas e URLs:
- gRPC: porta 50051
  - utilize o evans: evans -r repl
- GraphQL: [porta 8080 na URL](http://localhost:8080/)
- REST-API: porta 8000
  - utilize o vscode com plugin REST Client com os arquivos .http

## Makefile
Arquivo para automatizar operações de desenvolvimento, contem alguns targets:
- build: executa o build do projeto criando um binário chamado ordersystem
- run: executa o projeto (depende de services)
- test: executa os testes do projeto
- clean: limpa o ambiente e remove o binario do build
- services: coloca no ar o mysql e rabbitmq, aguarda o rabbitmq iniciar completamente
- services-shutdown: derruba os serviços de mysql e rabbitmq
- services-wipe: elimina a pasta .docker com os arquivos do volume do banco, solicita a senha para o sudo (depende de services-shutdown)
- generate-graphql: faz a geracao dos fontes

## Docker compose

O arquivo docker-compose.yaml define os servicos de mysql e rabbitmq para a aplicacao.
MySQL:
- nome do banco de dados: orders
- usuario: root
- senha: root
- porta: 3306

RabbitMQ:
- usurio: guest
- senha: guest
- postas:
  - 5672: amqp
  - 15672: geranciamento web

## Migrations

Para gerenciar as migrações de banco de dados é utilizado o https://github.com/golang-migrate/migrate
- As migrations estão no diretório migrations da raiz do projeto.
- Para criar as migrations foi utilizado o container docker da seguinte forma:
  ```docker run -v $(pwd)/migrations:/migrations --user $(id -u):$(id -g) migrate/migrate create -dir /migrations/ -ext sql orders```
- A aplicação da migration é feita pela aplicação no momento de iniciar.
- Os testes tambem aplicam as migrations.

## GRPC

Em caso de alteração dos arquivos proto (internal/infra/grpc/protofiles/order.proto):
- Instalar o protoc conforme documentação
- instalar os plugins:
  - go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  - go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
- executar o protoc para gerar os fontes do gRPC
  - diretamente: protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto
  - via make: make generate-grpc

## GraphQL

Em caso de alteração dos arquivos do GraphQL (internal/infra/graph/schema.graphqls)
Gerar os fontes para o GraphQL:
- diretamente: go run github.com/99designs/gqlgen generate
- via make: make generate-graphql

O arquivo gqlgen.yml configura a geracao de fontes do GraphQL

## DI Wire

Em caso de alteração do arquivo de definio do DI (cmd/ordersystem/wire.go):
- Instalar com:
  - go install github.com/google/wire/cmd/wire@latest
- Gerar os fontes de DI:
  - diretamente:
    - cd cmd/ordersystem
    - wire
  - via make: make wire-di

## Configuracao

A configuracao est no arquivo .env contendo os seguintes valores default:
- DB_DRIVER=mysql
- DB_HOST=localhost
- DB_PORT=3306
- DB_USER=root
- DB_PASSWORD=root
- DB_NAME=orders
- WEB_SERVER_PORT=:8000
- GRPC_SERVER_PORT=50051
- GRAPHQL_SERVER_PORT=8080


