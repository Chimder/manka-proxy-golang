FROM golang:1.23.3-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN apk --no-cache add ca-certificates

RUN go build -o ./example-golang ./cmd/api/main.go


FROM alpine:latest AS runner

WORKDIR /app
COPY --from=builder /app/example-golang .

ENV REDIS_URL=$REDIS_URL
ENV DB_URL=$DB_URL

EXPOSE 4000
ENTRYPOINT ["./example-golang"]