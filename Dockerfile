FROM golang:1.24.3-alpine AS builder
WORKDIR /app

# Copy files from host
# Copies from first path in host and pastes at second path (workdir) in container
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o Go_CRUD


FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /app/Go_CRUD .

CMD ["./Go_CRUD"]



