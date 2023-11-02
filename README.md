# Platform API
The entrypoint into the event flows of the microservice.
It is in charge of exposing REST API Endpoints, data sent to these endpoints are in turn published to the MQ.

# Dependencies
This project requires the following dependencies:
- Docker
- Kubernetes
- Helm
- RabbitMQ
- MongoDB

# Installation
This section covers how you can deploy this service for production


# Development
This section covers how you can deploy this service locally for development

## Cluster
```bash
# creates a new minikube cluster
# this assumes you're using minikube with a docker engine
minikube start
```

## Dependencies
To get the docker image for the project, you can build it locally:
```bash
# build with the dev file
docker build -t <image> -f ./docker/Dockerfile.dev .
```

