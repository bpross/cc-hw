FROM golang

ENV GO111MODULE=on

WORKDIR /cc

COPY go.mod .
COPY go.sum .
COPY . .

RUN go get -u golang.org/x/lint/golint
RUN go get github.com/onsi/ginkgo/ginkgo
