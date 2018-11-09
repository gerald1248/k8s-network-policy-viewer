# check for prerequisites
EXECUTABLES = docker jq
K := $(foreach exec,$(EXECUTABLES),\
	$(if $(shell which $(exec)),"",$(error "$(exec) not in path")))

# fetch name and namespace
NAME = $(shell jq -r .name values.yaml)
NAMESPACE = $(shell jq -r .namespace values.yaml)
GROUP = $(shell jq -r .group values.yaml)

# targets
build:
	make -C doc/
	docker build -t $(NAME) .
push:
	docker tag $(NAME) $(GROUP)/$(NAME):latest
	docker push $(GROUP)/$(NAME):latest
install:
	helm install --namespace=$(NAMESPACE) --name=$(NAME) .
delete:
	helm delete --purge $(NAME)
test:
	./Dockerfile_test
