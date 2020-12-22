FROM golang:1.14-alpine as builder

RUN apk update && apk add gcc libc-dev curl git

WORKDIR /app
COPY .env .env
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-linkmode=external' -a -installsuffix cgo -o main .

FROM alpine:3.11
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
EXPOSE 80