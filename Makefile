.PHONY: build-census build-postgres-dbinit \
minikube minikube-clean minikube-start minikube-stop minikube-tunnel \
postgres postgres-clean postgres-reset \
minikube-build-postgres-dbinit \
postgres-config postgres-storage postgres-deployment postgres-service \
postgres-service-clean postgres-deployment-clean postgres-storage-clean postgres-config-clean \
census census-clean census-rebuild \
minikube-build-census \
census-config census-deployment census-service \
census-config-clean census-deployment-clean census-service-clean


BUILD_VERSION ?= dev

#
# Local image builds
#

build-census:
	@docker build -t census:$(BUILD_VERSION) .

build-postgres-dbinit:
	@docker build -t postgres-dbinit:$(BUILD_VERSION) db


#
# Minikube commands
#

minikube: minikube-start postgres census

minikube-clean: census-clean postgres-clean minikube-stop

minikube-start:
	@minikube start

minikube-stop:
	@minikube stop

# Start minikube tunnel, which allows for connecting to the census service on localhost from the host machine
minikube-tunnel:
	@minikube tunnel


#
# Sample-microservice
#

census: minikube-build-census census-config census-deployment census-service
census-clean: census-service-clean census-deployment-clean census-config-clean
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

census-config:
	@kubectl create -f $(CENSUS_CONFIG)

census-config-clean:
	@kubectl delete -f $(CENSUS_CONFIG) --ignore-not-found

census-deployment:
	@kubectl create -f $(CENSUS_DEPLOYMENT)

census-deployment-clean:
	@kubectl delete -f $(CENSUS_DEPLOYMENT) --ignore-not-found

census-service:
	@kubectl create -f $(CENSUS_SERVICE)

census-service-clean:
	@kubectl delete -f $(CENSUS_SERVICE) --ignore-not-found


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
		docker build -t postgres-dbinit:$(BUILD_VERSION) db; \
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
	@kubectl delete -f $(POSTGRES_CLEANER_DEPLOYMENT) --ignore-not-found
	@kubectl create -f $(POSTGRES_CLEANER_DEPLOYMENT)
	@kubectl wait --for condition=complete -f $(POSTGRES_CLEANER_DEPLOYMENT)
