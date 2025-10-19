FROM golang:1.25-bookworm AS builder 

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/ws-gateway ./cmd/ws_gateway

FROM alpine:latest AS runner

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=builder /go/bin/ws-gateway .

RUN chmod +x ./ws-gateway

EXPOSE 8010 

CMD ["./ws-gateway"]