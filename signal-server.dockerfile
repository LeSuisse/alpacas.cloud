FROM golang:1.16.4-buster AS builder-go

WORKDIR /go/src/app
COPY . .

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get -y install gcc
RUN go build -v \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -modcacherw \
    -ldflags "-s -w -extldflags -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now" \
    -o signal-server cmd/alpacascloud-signal/main.go

FROM debian:buster-20210511-slim

ENV SHA256_LIB 4693cdfc8f49f4c7b23495a7330dbe2f024efebc95e7571f4331ac3e85765698
COPY --from=builder-go /go/src/app/signal-server /
ADD https://github.com/nanu-c/zkgroup/releases/download/v0.8.1/libzkgroup-linux_x86_64-v0.8.1.so /usr/lib/libzkgroup_linux_amd64.so
RUN echo "$SHA256_LIB /usr/lib/libzkgroup_linux_amd64.so" | sha256sum -c - && chmod 0444 /usr/lib/libzkgroup_linux_amd64.so

USER nobody
ENV STORAGE_DIRECTORY /storage
VOLUME /storage

CMD ["/signal-server"]
