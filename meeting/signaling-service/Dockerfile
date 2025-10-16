FROM golang:1.24-bookworm AS builder 

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/signaling-service ./cmd/signaling-service

FROM alpine:latest AS runner

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=builder /go/bin/signaling-service .

RUN chmod +x ./signaling-service

EXPOSE 8010 

CMD ["./signaling-service"]