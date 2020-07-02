FROM golang:1.14.4-buster AS builder-go

RUN apt-get update -y && apt-get install -y libpng-dev

WORKDIR /go/src/app
COPY . .

RUN go build -v \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -modcacherw \
    -ldflags "-s -w -extldflags -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now" \
    -o image-api-server cmd/alpacascloud/main.go

FROM node:14.5.0-buster-slim AS builder-web

COPY web/ /web/
WORKDIR /web/
RUN npm install && npm run build

FROM gcr.io/distroless/cc-debian10

COPY --from=builder-go /go/src/app/image-api-server /
COPY --from=builder-go /usr/lib/x86_64-linux-gnu/libpng16.so.16 /usr/lib/x86_64-linux-gnu/
COPY --from=builder-go /usr/lib/x86_64-linux-gnu/libpng16.so.16.36.0 /usr/lib/x86_64-linux-gnu/
COPY --from=builder-go /lib/x86_64-linux-gnu/libz.so.1 /lib/x86_64-linux-gnu/
COPY --from=builder-go /lib/x86_64-linux-gnu/libz.so.1.2.11 /lib/x86_64-linux-gnu/
COPY --from=builder-web /web/dist/ /web/dist/

USER nobody
EXPOSE 8080
ENV IMAGES_PATH /img

CMD ["/image-api-server"]
