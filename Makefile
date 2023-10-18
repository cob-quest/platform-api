IMAGE_TAG=registry.gitlab.com/cs302-2023/g3-team8/project/platform-api/main:latest
.DEFAULT_GOAL := run

# running
run:
	go run main.go

# linting
lint:
	golangci-lint run --enable-all

test:
	go test -v

# building
build:
	go build -o main main.go 

# clean
clean:
	rm ./main

# builds a docker image
docker:
	docker build -t ${IMAGE_TAG} -f ./docker/Dockerfile .

# builds a docker image and loads into minikube
minikube: docker
	minikube image rm ${IMAGE_TAG}
	minikube image load ${IMAGE_TAG}

