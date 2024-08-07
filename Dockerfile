FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN apk update && apk add --no-cache gcc musl-dev

RUN CGO_ENABLED=1 GOOS=linux go build -o app .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/.env .

COPY --from=builder /app/app .

COPY --from=builder /app/www ./www

EXPOSE 8080

CMD [ "./app" ]