FROM golang:1.17.8-alpine AS builder
RUN apk add --no-cache make gcc musl-dev linux-headers git ca-certificates
WORKDIR /app
COPY . .
RUN go build -ldflags '-w -s' .

FROM alpine:3.15.0
COPY --from=builder /app/metis-proxy /usr/local/bin/metis-proxy 
EXPOSE 8545
ENTRYPOINT [ "metis-proxy" ]
