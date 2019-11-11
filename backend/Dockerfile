FROM golang:1.13.4 AS builder
RUN go get github.com/gorilla/mux && go get github.com/lib/pq && go get github.com/nats-io/nats.go
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .

FROM scratch
COPY --from=builder /main ./
EXPOSE 8090
CMD ["/main"]