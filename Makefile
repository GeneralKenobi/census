.PHONY: build-census build-postgres-dbinit build-mongo-dbinit build-mongo-key-generator test e2e \
minikube minikube-clean minikube-start minikube-stop minikube-tunnel \
postgres postgres-clean postgres-reset \
minikube-build-postgres-dbinit \
postgres-config postgres-storage postgres-deployment postgres-service \
postgres-service-clean postgres-deployment-clean postgres-storage-clean postgres-config-clean \
mongo mongo-clean mongo-reset \
minikube-build-mongo-dbinit minikube-build-mongo-key-generator \
mongo-config mongo-storage mongo-deployment mongo-replica-set-init mongo-service \
mongo-service-clean mongo-deployment-clean mongo-storage-clean mongo-config-clean \
census census-clean census-rebuild \
minikube-build-census \
census-config census-tls-certs census-deployment census-service \
census-config-clean census-tls-certs-clean census-deployment-clean census-service-clean \
enable-postgres update-config-to-postgres enable-mongo update-config-to-mongo


BUILD_VERSION ?= dev

#
# Local image builds
#

build-census:
	@docker build -t census:$(BUILD_VERSION) .

build-postgres-dbinit:
	@docker build -t postgres-dbinit:$(BUILD_VERSION) db/postgres/init

build-mongo-dbinit:
	@docker build -t mongo-dbinit:$(BUILD_VERSION) db/mongo/init

build-mongo-key-generator:
	@docker build -t mongo-key-generator:$(BUILD_VERSION) db/mongo/keys



#
# Tests
#

test:
	@go test ./cmd/... ./pkg/... ./internal/...

CENSUS_HOST ?= localhost
CENSUS_PORT ?= 8443
CENSUS_PROTOCOL ?= https
e2e:
	@# E2e test results can't be cached because they don't depend directly on the tested code.
	@go clean -testcache ./test/e2e/...
	@CENSUS_HOST=$(CENSUS_HOST) \
	CENSUS_PORT=$(CENSUS_PORT) \
	CENSUS_PROTOCOL=$(CENSUS_PROTOCOL) \
	go test -v ./test/e2e/...


#
# Minikube commands
#

minikube: minikube-start

minikube-clean: census-clean postgres-clean mongo-clean minikube-stop

minikube-start:
	@minikube start

minikube-stop:
	@minikube stop

# Start minikube tunnel, which allows for connecting to the census service on localhost from the host machine
minikube-tunnel:
	@minikube tunnel


#
# Census
#

census: minikube-build-census census-tls-certs census-config census-deployment census-service
census-clean: census-service-clean census-deployment-clean census-config-clean census-tls-certs-clean
# Rebuild and redeploy census, update configuration
census-rebuild: minikube-build-census census-deployment-clean census-config-clean census-config census-deployment

minikube-build-census:
	@{ \
		eval $$(minikube docker-env); \
		docker build -t census:$(BUILD_VERSION) .; \
	}

CENSUS_YAMLS_DIR=deployments/kubernetes/local/census
CENSUS_CONFIG=$(CENSUS_YAMLS_DIR)/config.yaml
CENSUS_DEPLOYMENT=$(CENSUS_YAMLS_DIR)/deployment.yaml
CENSUS_SERVICE=$(CENSUS_YAMLS_DIR)/service.yaml
CENSUS_TLS_SECRET_NAME=census-cert

census-config:
	@kubectl create -f $(CENSUS_CONFIG)

census-config-clean:
	@kubectl delete -f $(CENSUS_CONFIG) --ignore-not-found

census-tls-certs:
	@openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes -subj "/CN=census.szyzog.com" -keyout certificates/census.key -out certificates/census.crt
	@kubectl create secret tls $(CENSUS_TLS_SECRET_NAME) --cert certificates/census.crt --key certificates/census.key
	@kubectl label secret $(CENSUS_TLS_SECRET_NAME) app=census

census-tls-certs-clean:
	@kubectl delete secret $(CENSUS_TLS_SECRET_NAME) --ignore-not-found
	@rm -f certificates/census.crt certificates/census.key

census-deployment:
	@kubectl create -f $(CENSUS_DEPLOYMENT)

census-deployment-clean:
	@kubectl delete -f $(CENSUS_DEPLOYMENT) --ignore-not-found

census-service:
	@kubectl create -f $(CENSUS_SERVICE)

census-service-clean:
	@kubectl delete -f $(CENSUS_SERVICE) --ignore-not-found


#
# Switching between DBs
#

enable-postgres: census-deployment-clean census-config-clean update-config-to-postgres census-config census-deployment

update-config-to-postgres:
	@sed -i '' -e 's/"database": "mongo"/"database": "postgres"/g' deployments/kubernetes/local/census/config.yaml

enable-mongo: census-deployment-clean census-config-clean update-config-to-mongo census-config census-deployment

update-config-to-mongo:
	@sed -i '' -e 's/"database": "postgres"/"database": "mongo"/g' deployments/kubernetes/local/census/config.yaml


#
# Postgres
#

postgres: minikube-build-postgres-dbinit postgres-config postgres-storage postgres-deployment postgres-service
postgres-clean: postgres-service-clean postgres-deployment-clean postgres-storage-clean postgres-config-clean
# Gracefully delete DB content and reinitialize the schema
postgres-reset: postgres-deployment-clean postgres-clear-volume minikube-build-postgres-dbinit postgres-deployment

minikube-build-postgres-dbinit:
	@{ \
		eval $$(minikube docker-env); \
		docker build -t postgres-dbinit:$(BUILD_VERSION) db/postgres/init; \
	}

POSTGRES_YAMLS_DIR=deployments/kubernetes/local/postgres
POSTGRES_CONFIG=$(POSTGRES_YAMLS_DIR)/config.yaml
POSTGRES_DEPLOYMENT=$(POSTGRES_YAMLS_DIR)/deployment.yaml
POSTGRES_SERVICE=$(POSTGRES_YAMLS_DIR)/service.yaml
POSTGRES_STORAGE=$(POSTGRES_YAMLS_DIR)/storage.yaml
POSTGRES_CLEANER_DEPLOYMENT=$(POSTGRES_YAMLS_DIR)/cleaner-deployment.yaml

postgres-config:
	@kubectl create -f $(POSTGRES_CONFIG)

postgres-config-clean:
	@kubectl delete -f $(POSTGRES_CONFIG) --ignore-not-found

postgres-storage:
	@kubectl create -f $(POSTGRES_STORAGE)

postgres-storage-clean:
	@kubectl delete -f $(POSTGRES_STORAGE) --ignore-not-found

postgres-deployment:
	@kubectl create -f $(POSTGRES_DEPLOYMENT)

postgres-deployment-clean:
	@kubectl delete -f $(POSTGRES_DEPLOYMENT) --ignore-not-found

postgres-service:
	@kubectl create -f $(POSTGRES_SERVICE)

postgres-service-clean:
	@kubectl delete -f $(POSTGRES_SERVICE) --ignore-not-found

# Delete DB content
postgres-clear-volume:
	@kubectl wait --for delete pod --selector=app=postgres
	@kubectl delete -f $(POSTGRES_CLEANER_DEPLOYMENT) --ignore-not-found
	@kubectl create -f $(POSTGRES_CLEANER_DEPLOYMENT)
	@kubectl wait --for condition=complete -f $(POSTGRES_CLEANER_DEPLOYMENT)


#
# Mongo
#

mongo: minikube-build-mongo-dbinit minikube-build-mongo-key-generator mongo-config mongo-storage mongo-deployment mongo-replica-set-init mongo-service
mongo-clean: mongo-service-clean mongo-deployment-clean mongo-storage-clean mongo-config-clean
# Gracefully delete DB content and reinitialize the database
mongo-reset: mongo-deployment-clean mongo-clear-volume minikube-build-mongo-dbinit minikube-build-mongo-key-generator mongo-deployment mongo-replica-set-init

minikube-build-mongo-dbinit:
	@{ \
		eval $$(minikube docker-env); \
		docker build -t mongo-dbinit:$(BUILD_VERSION) db/mongo/init; \
	}

minikube-build-mongo-key-generator:
	@{ \
		eval $$(minikube docker-env); \
		docker build -t mongo-key-generator:$(BUILD_VERSION) db/mongo/keys; \
	}


MONGO_YAMLS_DIR=deployments/kubernetes/local/mongo
MONGO_CONFIG=$(MONGO_YAMLS_DIR)/config.yaml
MONGO_DEPLOYMENT=$(MONGO_YAMLS_DIR)/deployment.yaml
MONGO_SERVICE=$(MONGO_YAMLS_DIR)/service.yaml
MONGO_STORAGE=$(MONGO_YAMLS_DIR)/storage.yaml
MONGO_CLEANER_DEPLOYMENT=$(MONGO_YAMLS_DIR)/cleaner-deployment.yaml

mongo-config:
	@kubectl create -f $(MONGO_CONFIG)

mongo-config-clean:
	@kubectl delete -f $(MONGO_CONFIG) --ignore-not-found

mongo-storage:
	@kubectl create -f $(MONGO_STORAGE)

mongo-storage-clean:
	@kubectl delete -f $(MONGO_STORAGE) --ignore-not-found

mongo-deployment:
	@kubectl create -f $(MONGO_DEPLOYMENT)

MONGO_INIT_REPLICA_SET_JS=if (rs.status().code == 94) {print(\"Initializing replica set\"); rs.initiate( {_id : \"rs0\", members: [{ _id: 0, host: \"mongo:27017\" }]});} else { print(\"Replica set already initialized\")}
mongo-replica-set-init:
	@kubectl wait pods --for condition=Ready --selector app=mongo
	@# Something's wrong with the k8s network and when mongo tries to connect to itself via the 'mongo' k8s service it fails.
	@# To fix this we map the 'mongo' service DNS in /etc/hosts to localhost to bypass the k8s network (there's only 1 mongo pod anyway).
	@kubectl exec -it -c mongo deploy/mongo -- /bin/bash -c "echo '127.0.0.1 mongo' >> /etc/hosts"
	@kubectl exec -it -c mongo deploy/mongo -- /bin/bash -c "mongo -u \$${MONGO_INITDB_ROOT_USERNAME} -p \$${MONGO_INITDB_ROOT_PASSWORD} --eval '$(MONGO_INIT_REPLICA_SET_JS)'"

mongo-deployment-clean:
	@kubectl delete -f $(MONGO_DEPLOYMENT) --ignore-not-found

mongo-service:
	@kubectl create -f $(MONGO_SERVICE)

mongo-service-clean:
	@kubectl delete -f $(MONGO_SERVICE) --ignore-not-found

# Delete DB content
mongo-clear-volume:
	@kubectl wait --for delete pod --selector=app=mongo
	@kubectl delete -f $(MONGO_CLEANER_DEPLOYMENT) --ignore-not-found
	@kubectl create -f $(MONGO_CLEANER_DEPLOYMENT)
	@kubectl wait --for condition=complete -f $(MONGO_CLEANER_DEPLOYMENT)
