FROM golang:1.19.2

WORKDIR /usr/src/app

#### Cache Dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify
####

COPY . .

RUN go build -o main ./cmd/main.go

CMD ["./main"]