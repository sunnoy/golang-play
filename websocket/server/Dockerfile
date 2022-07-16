##
## Build
##
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /server

##
## Deploy
##
FROM centos:7

WORKDIR /

COPY --from=build /server /server

EXPOSE 2022

ENTRYPOINT ["/server"]