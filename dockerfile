FROM golang:1.23.11-alpine as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /app/stress cmd/main.go

# FROM scratch
FROM alpine:latest

COPY --from=builder /app .

ENTRYPOINT ["/stress"]