FROM golang:latest

WORKDIR /var/api

COPY . ./

RUN go mod download
RUN go build ./cmd/api
RUN mkdir -p /var/api/images

CMD ["./api"]