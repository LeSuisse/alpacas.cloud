FROM node:19.9.0-alpine3.16 AS builder-web

COPY cmd/alpacascloud/web/ /web/
WORKDIR /web/
RUN npm install && npm run build

FROM golang:1.20.3-alpine3.16 AS builder-go

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

FROM cgr.dev/chainguard/alpine-base:latest

RUN apk add --no-cache --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community vips

COPY --from=builder-go /go/src/app/image-api-server /

USER nobody
EXPOSE 8080
ENV IMAGES_PATH /img

CMD ["/image-api-server"]
