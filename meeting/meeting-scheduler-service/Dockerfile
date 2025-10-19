FROM golang:1.25-bookworm AS builder 

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/meeting-scheduler ./cmd/meeting-scheduler-service

FROM alpine:latest AS runner

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=builder /go/bin/meeting-scheduler .

RUN chmod +x ./meeting-scheduler

EXPOSE 8003 

CMD ["./meeting-scheduler"]