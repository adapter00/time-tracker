FROM golang:1.21.1
ENV CGO_ENABLED 1
RUN mkdir /go/src/app
WORKDIR /go/src/app
ADD . /go/src/app
RUN go mod download
RUN go install github.com/cespare/reflex@latest
CMD cd /go/src/app && reflex -r '(\.go$|go\.mod)' -s go run .
