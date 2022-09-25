FROM node:18.9.0-alpine3.16 AS builder-web

COPY cmd/alpacascloud/web/ /web/
WORKDIR /web/
RUN npm install && npm run build

FROM golang:1.19.1-alpine3.16 AS builder-go

RUN apk add --no-cache vips-dev gcc libc-dev pkgconfig

WORKDIR /go/src/app
COPY . .
COPY --from=builder-web /web/dist/ ./cmd/alpacascloud/web/dist/

RUN go build -v \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -modcacherw \
    -ldflags "-s -w -extldflags -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now" \
    -o image-api-server cmd/alpacascloud/main.go

FROM cgr.dev/chainguard/alpine-base:latest@sha256:9967caa74bbb6b11402a963735f1ae202af362ecf280c169a26865df862cf2eb

RUN apk add --no-cache --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community vips

COPY --from=builder-go /go/src/app/image-api-server /

USER nobody
EXPOSE 8080
ENV IMAGES_PATH /img

CMD ["/image-api-server"]
