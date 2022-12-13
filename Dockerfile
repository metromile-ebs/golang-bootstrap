FROM golang:1.19.2

WORKDIR /usr/src/app

#### Cache Dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify
####

COPY . .

RUN wget https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem

RUN go build -o main ./cmd/main.go

CMD ["./main"]