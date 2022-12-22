FROM golang:1.19.2

WORKDIR /usr/src/app

ARG USERNAME
ENV USERNAME=$USERNAME
ARG TOKEN
ENV TOKEN=$TOKEN

RUN go env -w GOPRIVATE=github.com/metromile-ebs
RUN git config \
    --global \
    url."https://${USERNAME}:${TOKEN}@github.com".insteadOf \
    "https://github.com"

#### Cache Dependencies
COPY go.mod go.sum ./
RUN GIT_TERMINAL_PROMPT=1 go mod download && go mod verify
####

COPY . .

RUN go build -o main ./cmd/main.go

CMD ["./main"]
