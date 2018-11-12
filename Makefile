# check for prerequisites
EXECUTABLES = docker jq
K := $(foreach exec,$(EXECUTABLES),\
	$(if $(shell which $(exec)),"",$(error "$(exec) not in path")))

# fetch name and namespace
NAME = $(shell jq -r .name chart/values.yaml)
NAMESPACE = $(shell jq -r .namespace chart/values.yaml)
GROUP = $(shell jq -r .group chart/values.yaml)

# targets
build:
	docker build -t $(NAME) .
push:
	docker tag $(NAME) $(GROUP)/$(NAME):latest
	docker push $(GROUP)/$(NAME):latest
test:
	./Dockerfile_test
