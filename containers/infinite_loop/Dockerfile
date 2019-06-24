FROM golang:alpine
WORKDIR /go/src/github.com/step/uruk/containers/infinite_loop
ADD . .
RUN CGO_ENABLED=0 go build -o ./runner ./main.go

FROM scratch
WORKDIR /app
COPY --from=0 /go/src/github.com/step/uruk/containers/infinite_loop/runner ./runner
ENTRYPOINT ["/app/runner"]
