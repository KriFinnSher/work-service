FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /work-service .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /work-service .

ENV GRPC_PORT=8082
ENV HTTP_PORT=8083

EXPOSE 8082 8083

CMD ["./work-service"]