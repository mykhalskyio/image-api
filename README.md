![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mykhalskyio/image-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/mykhalskyio/image-api)](https://goreportcard.com/report/github.com/mykhalskyio/image-api)
![Github Repository Size](https://img.shields.io/github/repo-size/mykhalskyio/image-api)

# Image API

## Requirements
* Git
* Docker

## How to run with make

Clone repository

    $ git clone https://github.com/mykhalskyio/image-api.git

Run make file
    
    $ make build
    
Run migrate

    $ make migrate-init

## How to run with docker

Clone repository

    $ git clone https://github.com/mykhalskyio/image-api.git

Run docker-compose

    $ docker-compose up --build api

Run migration

    $ docker run -v ./schema:/migrations --network host migrate/migrate     -path=/migrations/ -database postgres://postgres:qwerty123@localhost:54320/postgres?sslmode=disable up

## Endpoints
* POST   - /api/upload (form-data, image)
* GET    - /api/download/id?quality=100/75/50/25
* DELETE - /api/delete/id