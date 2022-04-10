FROM golang:1.17 as builder

WORKDIR /mailsender
COPY . /mailsender

RUN export CGO_ENABLED=0 && go build -o mailsender ./

FROM alpine:3.15
RUN apk update && apk add --no-cache bash
COPY --from=builder /mailsender /mailsender

VOLUME "/var/lib/mailsender"

ENV PORT 8082

EXPOSE ${PORT}

WORKDIR /mailsender

ENTRYPOINT ["./mailsender", "start"]