FROM golang:1.17.8-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags '-w -s' .

FROM alpine:3.15.0
COPY --from=builder /app/metis-proxy /usr/local/bin/metis-proxy 
EXPOSE 8545
ENTRYPOINT [ "metis-proxy" ]
