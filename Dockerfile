FROM golang:1.22

WORKDIR /usr/src/music

COPY go.mod ./

# COPY go.mod go.sum ./

RUN go mod tidy

RUN go mod download && go mod verify

COPY . .

# install migrate cli
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate

RUN go build -v -o /usr/local/bin/music/ ./...
