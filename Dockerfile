FROM golang:1.24-alpine

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "cmd/main.go"]