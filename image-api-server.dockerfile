FROM node:24.4.1-alpine3.21 AS builder-web

COPY cmd/alpacascloud/web/ /web/
WORKDIR /web/
RUN npm install && npm run build

FROM golang:1.24-alpine3.21 AS builder-go

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

FROM alpine:3.22

RUN apk add --no-cache vips

COPY --from=builder-go /go/src/app/image-api-server /

USER nobody
EXPOSE 8080
ENV IMAGES_PATH /img

CMD ["/image-api-server"]
