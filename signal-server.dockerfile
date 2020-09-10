FROM golang:1.15.2-buster AS builder-go

WORKDIR /go/src/app
COPY . .

RUN go build -v \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -modcacherw \
    -ldflags "-s -w -extldflags -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now" \
    -o signal-server cmd/alpacascloud-signal/main.go

FROM gcr.io/distroless/base

COPY --from=builder-go /go/src/app/signal-server /

USER nobody
ENV STORAGE_DIRECTORY /storage
VOLUME /storage

CMD ["/signal-server"]
