FROM golang:1.16.2-alpine3.12 AS builder-go

WORKDIR /go/src/app
COPY . .

RUN apk add --no-cache gcc musl-dev
RUN go build -v \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -modcacherw \
    -ldflags "-s -w -extldflags -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now" \
    -o signal-server cmd/alpacascloud-signal/main.go

FROM alpine:3.13.4

COPY --from=builder-go /go/src/app/signal-server /

USER nobody
ENV STORAGE_DIRECTORY /storage
VOLUME /storage

CMD ["/signal-server"]
