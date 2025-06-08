FROM golang:1.24-alpine

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app/backend

COPY ../backend /app/backend

RUN go mod download

CMD ["go", "run", "main.go"]
