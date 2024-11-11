FROM golang:1.23.3-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN apk --no-cache add ca-certificates

RUN go build -o ./example-golang ./main.go

FROM alpine:latest AS runner

WORKDIR /app
COPY --from=builder /app/example-golang .

EXPOSE 8080
ENTRYPOINT ["./example-golang"]