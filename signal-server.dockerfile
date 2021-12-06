FROM golang:1.17.4-buster AS builder-go

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

ENV SHA256_LIB 6e743a80b1334e655882286fae047ab701094594166aeea88fc4089f1fdb8f57
ADD https://github.com/nanu-c/zkgroup/releases/download/v0.8.8/libzkgroup_linux_x86_64-v0.8.8.so /usr/lib/libzkgroup_linux_x86_64.so
RUN echo "$SHA256_LIB /usr/lib/libzkgroup_linux_x86_64.so" | sha256sum -c - && chmod 0444 /usr/lib/libzkgroup_linux_x86_64.so

FROM gcr.io/distroless/cc-debian10

COPY --from=builder-go /go/src/app/signal-server /
COPY --from=builder-go /usr/lib/libzkgroup_linux_x86_64.so /usr/lib/libzkgroup_linux_x86_64.so

USER nobody
ENV STORAGE_DIRECTORY /storage
VOLUME /storage

CMD ["/signal-server"]
