FROM golang:1.17.2-buster
ENV CGO_ENABLED 0
RUN mkdir /go/src/app
WORKDIR /go/src/app
ADD . /go/src/app
RUN go mod download
RUN go get -u github.com/cespare/reflex
CMD cd /go/src/app && reflex -r '(\.go$|go\.mod)' -s go run .
