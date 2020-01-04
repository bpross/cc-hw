FROM golang

ENV GO111MODULE=on

WORKDIR /

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/server/main.go

EXPOSE 8080
ENTRYPOINT ["/main"]
