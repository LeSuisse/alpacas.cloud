FROM golang:1.14.3-buster AS builder-go

WORKDIR /go/src/app
COPY . .

RUN go build -v -o signal-server cmd/alpacascloud-signal/main.go

FROM gcr.io/distroless/base

COPY --from=builder-go /go/src/app/signal-server /

USER nobody
ENV STORAGE_DIRECTORY /storage
VOLUME /storage

CMD ["/signal-server"]
