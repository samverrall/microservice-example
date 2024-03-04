# Golang DDD / gRPC / Hexagonal Architecture Example

This repository holds an example of how you may structure a Golang microservice. The structure of this service follows principles 
from Hexagonal Architecture (ports/adapters) and Domain Driven Design. It exposes a gRPC transport layer (adapter) for internal communication 
between other services. The idea of this repository is to share an approach to how you could structure microservices in Golang.

## Project structure 

`internal/app/` holds the core application layer as "services".

## Running the service locally 

```sh 
make build 
make run 
```

## Generating Proto 

```sh
make proto-health
```

