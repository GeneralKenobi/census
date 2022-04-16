# Census

Go service for managing people.

## Requirements

This software is developed and tested on MacOS.

- Go 1.18
- Docker
- Minikube
- CLI binaries: make, rm, openssl

## Quickstart

#### Build census image

```shell
make
```

#### Run unit tests

```shell
make test
```

#### Run e2e tests

```shell
# By default expects census to listen on localhost:8443 with https
make e2e

# Set custom census deployment information with environment variables
make e2e CENSUS_PROTOCOL=http CENSUS_HOST=53.142.61.11 CENSUS_PORT=80
```

#### Start minikube

```shell
minikube start
# Or
make minikube
```

#### Start minikube tunnelling to access cluster services from local machine

```shell
minikube tunnel
# Or 
make minikube-tunnel
```

#### Stop minikube and clean up everything

```shell
make minikube-clean
```

#### Build and deploy census with configuration, services, etc.

```shell
make census
```

#### Rebuild and redeploy census, and update the configuration

```shell
make census-rebuild
```

#### Clean census and its configuration, services, etc.

```shell
make census-clean
```

#### Deploy postgres DB with configuration, services, etc.

```shell
make postgres
```

#### Switch to postgres DB

```shell
make enable-postgres
```

#### Reset and redeploy postgres DB (start from scratch)

```shell
make postgres-reset
```

#### Clean postgres DB and its configuration, services, etc.

```shell
make postgres-clean
```

#### Deploy mongo DB with configuration, services, etc.

```shell
make mongo
```

#### Switch to mongo DB

```shell
make enable-mongo
```

#### Reset and redeploy mongo DB (start from scratch)

```shell
make mongo-reset
```

#### Clean mongo DB and its configuration, services, etc.

```shell
make mongo-clean
```

## Functional requirements

This application allows for creating, reading, updating, and deleting `person` entities.

## Sample requests

#### Create a person

```shell
curl https://localhost:8443/api/person -k -X POST -H "Content-Type: application/json" -d '{"name":"John","surname":"Smith","email":"john.smith@test.com","dateOfBirth":"1995-09-13T00:00:00Z","hobby":"Jogging"}'
HTTP200 {"id":"625aed07d292114997997af5"}
```

#### Get a person

```shell
curl https://localhost:8443/api/person/625aed07d292114997997af5 -k -X GET
HTTP200 {"id":"625aed07d292114997997af5","name":"John","surname":"Smith","email":"john.smith@test.com","dateOfBirth":"1995-09-13T00:00:00Z","hobby":"Jogging","createdAt":"2022-04-16T16:21:27.038Z","lastModifiedAt":"2022-04-16T16:21:27.038Z"}
```

#### Update a person

```shell
curl https://localhost:8443/api/person/625aed07d292114997997af5 -k -X PUT -H "Content-Type: application/json" -d '{"name":"James","surname":"Smith","dateOfBirth":"1995-09-13T00:00:00Z","hobby":"Jogging"}'
HTTP200
```

#### Delete a person

```shell
curl https://localhost:8443/api/person/625aed07d292114997997af5 -k -X DELETE
HTTP200
```

## Repository structure

- `api` - API schema
- `certificates` - Target directory for self-signed certificate generation
- `cmd/census` - Entrypoint with `main` function and flag definitions
- `db/mongo` - MongoDB setup
- `db/postgres` - PostgresDB setup
- `deployments/kubernetes/local` - Kubernetes resource files for local development on minikube
- `internal/api` - HTTP server code
- `internal/config` - Configuration types and configuration loading
- `internal/db` - Database interfaces
- `internal/db/model` - Entity structs
- `internal/db/mongo` - MongoDB integration
- `internal/db/postgres` - PostgresDB integration
- `internal/service` - Services implementing business logic
- `pkg/api/apimodel` - API DTOs
- `pkg/api/client` - API client library
- `pkg/mdctx` - Logging framework with support for request-aware logging (it enhances logs with information like
  correlation ID or request path)
- `pkg/shutdown` - Context supporting graceful shutdown of the application
- `test/e2e` - End-to-end tests


## How to connect to DBs

Deployments of both DBs include an additional service that exposes them via `minikube tunnel` on the host machine on localhost on the DB's port.
After running `make minikube-tunnel` you can access
- PostgresDB on `localhost:5432`
- MongoDB on `localhost:27017`
