FROM golang:1.17 as builder

WORKDIR /mailsender
COPY . /mailsender

RUN export CGO_ENABLED=0 && go build -o mailsender ./

FROM alpine:3.15
RUN apk update && apk add --no-cache bash
COPY --from=builder /mailsender /mailsender

VOLUME "/var/lib/mailsender"

EXPOSE 8082/tcp
WORKDIR /mailsender
ENTRYPOINT ["./mailsender", "start", "-p", "8082"]