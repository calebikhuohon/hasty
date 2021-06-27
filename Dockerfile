FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN mkdir -p /app/bin
RUN go build \
  -o bin \
  ./cmd/...

FROM alpine:latest
RUN apk --no-cache --update add ca-certificates
COPY --from=builder /app/bin/* /app/bin/
COPY --from=builder /app/internal/storage/migrations/* /app/migrations/
CMD ["/usr/bin/false"]