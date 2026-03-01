
FROM golang:1.25-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/api/main.go



FROM alpine:latest

WORKDIR /

COPY --from=builder /app/myapp /myapp

EXPOSE 8000

ENTRYPOINT ["/myapp"]