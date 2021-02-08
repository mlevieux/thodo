#
# NOTE: THIS DOCKERFILE IS GENERATED VIA "apply-templates.sh"
#
# PLEASE DO NOT EDIT IT DIRECTLY.
#

FROM golang:1.15.7
RUN go get github.com/mlevieux/thodo/...

WORKDIR /go/src/github.com/mlevieux/thodo

RUN go build -o thodo-back ./back/

EXPOSE 8080
CMD ["./thodo-back"]