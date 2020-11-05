FROM golang:1.15.3-alpine3.12 AS builder-go

RUN apk add --no-cache vips-dev gcc libc-dev pkgconfig

WORKDIR /go/src/app
COPY . .

RUN go build -v \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -modcacherw \
    -ldflags "-s -w -extldflags -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now" \
    -o image-api-server cmd/alpacascloud/main.go

FROM node:15.1.0-alpine3.12 AS builder-web

COPY web/ /web/
WORKDIR /web/
RUN npm install && npm run build

FROM alpine:3.12.1

RUN apk add --no-cache vips

COPY --from=builder-go /go/src/app/image-api-server /
COPY --from=builder-web /web/dist/ /web/dist/

USER nobody
EXPOSE 8080
ENV IMAGES_PATH /img

CMD ["/image-api-server"]
