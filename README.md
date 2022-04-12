# Census

Go service for managing people.

## Requirements

- Go 1.18
- Docker
- Minikube

## Quickstart

#### Build census image

```shell
make
```

#### Start minikube and setup everything

```shell
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

#### Reset and redeploy postgres DB (start from scratch)

```shell
make postgres-reset
```

#### Clean postgres DB and its configuration, services, etc.

```shell
make postgres-clean
```

## Functional requirements

TODO

## Sample requests

TODO
