FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
COPY ./config ./config 
EXPOSE 8081
ENTRYPOINT ["/app"]
CMD ["--config", "./config/config.yaml"]
LABEL Name=server Version=0.0.1
