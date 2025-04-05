FROM golang:1.24-alpine AS builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app examples/cluster/main.go

FROM scratch
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/app ./app
EXPOSE 8080
ENTRYPOINT ["/usr/src/app/app"]
