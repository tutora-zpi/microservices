FROM golang:1.25-bookworm AS builder 

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/recording-service ./cmd/recorder_service

FROM alpine:latest AS runner

RUN apk add --no-cache libc6-compat ffmpeg

WORKDIR /app

COPY --from=builder /go/bin/recording-service .

RUN chmod +x ./recording-service

EXPOSE 8050 

CMD ["./recording-service"]
