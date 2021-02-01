#
# NOTE: THIS DOCKERFILE IS GENERATED VIA "apply-templates.sh"
#
# PLEASE DO NOT EDIT IT DIRECTLY.
#

FROM golang:1.15.7
RUN go get github.com/gorilla/mux
RUN go get github.com/json-iterator/go

WORKDIR /go/src/
COPY src thodo/src

WORKDIR /go/src/thodo/src
RUN go build -o thodo .

EXPOSE 8080
CMD ["./thodo"]