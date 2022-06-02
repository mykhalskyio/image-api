# Image API

## Requirements
* Git
* Docker

## How to run

Clone repository

    $ git clone https://github.com/mykhalskyio/image-api.git

Run make file
    
    $ make build
    
Run migrate

    $ make migrate-init

## Endpoints
* /api/upload (form-data, image)
* /api/download/id?quality=100/75/50/25
* /api/delete/id