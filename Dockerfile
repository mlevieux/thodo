#
# NOTE: THIS DOCKERFILE IS GENERATED VIA "apply-templates.sh"
#
# PLEASE DO NOT EDIT IT DIRECTLY.
#

FROM mysql/mysql-server:latest

FROM golang:1.15.7
RUN go get github.com/gorilla/mux
RUN go get github.com/json-iterator/go
RUN go get github.com/go-sql-driver/mysql

WORKDIR /go/src/
COPY src github.com/mlevieux/thodo/src

WORKDIR /go/src/github.com/mlevieux/thodo/src
RUN go build -o thodo .

EXPOSE 8080
CMD ["./thodo"]