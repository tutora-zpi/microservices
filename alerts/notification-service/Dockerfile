FROM golang:1.24-bookworm AS builder 

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/notification-service ./cmd/notification-service

FROM alpine:latest AS runner

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=builder /go/bin/notification-service .

RUN chmod +x ./notification-service

EXPOSE 8888 

CMD ["./notification-service"]