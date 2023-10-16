#!/usr/bin/env bash
docker build -t registry.gitlab.com/cs302-2023/g3-team8/project/platform-api/main -f ./docker/Dockerfile .
minikube image load registry.gitlab.com/cs302-2023/g3-team8/project/platform-api/main
