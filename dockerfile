FROM golang:1.17.0-buster AS builder

WORKDIR /app
RUN apt-get update && apt-get install -y protobuf-compiler golang-goprotobuf-dev
COPY . /app/
RUN make build

FROM alpine:latest
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY --from=builder /app/bin/ova-purchase-api .
COPY --from=builder /app/.env .
CMD ["/app/ova-purchase-api"]
