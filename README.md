# dockerapp

sample go lang 3 layers app run on top of docker

### Prerequisite
- Docker version 18.03.1-ce or later
- Go 1.10.1 or later

### Build & Installation
```sh
$ go get -u
$ make build-image
$ make build-docker
$ docker-compose up
```

### Docker stack
- Ubuntu 18.04 + ca-certificates
- Postgres 10.4
- Nginx 1.13