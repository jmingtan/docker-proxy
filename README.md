# Docker Proxy

Creates a localhost proxy for Docker connections. This is useful when you are running Docker within a VM, but would like to access containers via localhost.

## Dependencies

    go get github.com/fsouza/go-dockerclient

## To use

    go run main.go

Requires the necessary environment variables for accessing the docker VM to be present (e.g. DOCKER_HOST, DOCKER_TLS_VERIFY, DOCKER_CERT_PATH). Usually this can be taken care of by a provisioning tool such as [docker-machine](https://docs.docker.com/machine/reference/env/).