FROM golang:alpine
WORKDIR /go/src/github.com/step/uruk/

ADD . .
RUN apk update && apk add --no-cache git ca-certificates make && go get ./... && make uruk

FROM golang:alpine
WORKDIR /app
COPY --from=0 /go/src/github.com/step/uruk/bin/uruk ./uruk
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
ENTRYPOINT [ "/app/uruk" ]
CMD ["-redis-address","172.17.0.2:6379"]