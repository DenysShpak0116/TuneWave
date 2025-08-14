FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /usr/bin/app
COPY ./config ./app/config
WORKDIR /app
EXPOSE 8081
ENTRYPOINT ["/usr/bin/app"]
CMD ["--config", "./config/config.yaml"]
LABEL Name=server Version=0.0.1
