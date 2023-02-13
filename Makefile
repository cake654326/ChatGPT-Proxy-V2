VERSION=$(shell git describe --always --match "v[0-9]*" HEAD)
DOCKER_BASE=acheong08
NAME=chatgpt-proxy-v2

.PHONY: build-docker
build-docker:
	docker buildx build --push --platform linux/amd64 -t $(DOCKER_BASE)/$(NAME):$(VERSION) .
	docker buildx build --push --platform linux/amd64 -t $(DOCKER_BASE)/$(NAME):latest .

.PHONY: run-docker
run-docker:
	docker stop $(NAME) ; docker rm $(NAME) ; docker pull $(DOCKER_BASE)/$(NAME) ;
	docker run -itd --name $(NAME) \
    -p 10101:10101 \
    $(DOCKER_BASE)/$(NAME)