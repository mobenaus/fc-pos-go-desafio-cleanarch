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
## Makefile
Arquivo para automatizar operações de desenvolvimento, contem alguns targets:
- build: executa o build do projeto criando um binário chamado ordersystem
- run: executa o projeto (depende de services)
- test: executa os testes do projeto
- clean: limpa o ambiente e remove o binario do build
- services: coloca no ar o mysql e rabbitmq, aguarda o rabbitmq iniciar completamente
- services-shutdown: derruba os serviços de mysql e rabbitmq
- services-wipe: elimina a pasta .docker com os arquivos do volume do banco (solicita o password para o sudo)

## Migrations

Para gerenciar as migrações de banco de dados é utilizado o https://github.com/golang-migrate/migrate
- As migrations estão no diretório migrations da raiz do projeto.
- Para criar as migrations foi utilizado o container docker da seguinte forma:
  ```docker run -v $(pwd)/migrations:/migrations --user $(id -u):$(id -g) migrate/migrate create -dir /migrations/ -ext sql orders```
- A aplicação da migration é feita pela aplicação no momento de iniciar.
- Os testes tambem aplicam as migrations.
