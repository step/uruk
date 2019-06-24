FROM golang:alpine
WORKDIR /go/src/github.com/step/Uruk/sample
ADD . .
RUN CGO_ENABLE=0 go build -o ./runner ./main.go

FROM scratch
WORKDIR /app
COPY --from=0 /go/src/github.com/step/Uruk/sample/runner ./runner
ENTRYPOINT ["/app/runner"]
