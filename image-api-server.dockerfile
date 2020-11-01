FROM golang:1.15.3-buster AS builder-go

RUN apt-get update -y && apt-get install --no-install-recommends -y libvips-dev

WORKDIR /go/src/app
COPY . .

RUN go build -v \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -modcacherw \
    -ldflags "-s -w -extldflags -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now" \
    -o image-api-server cmd/alpacascloud/main.go

FROM node:15.0.1-buster-slim AS builder-web

COPY web/ /web/
WORKDIR /web/
RUN npm install && npm run build

FROM ubuntu:focal-20201008

RUN apt-get update -y && apt-get install --no-install-recommends -y libvips && rm -rf /var/lib/apt/lists/*

COPY --from=builder-go /go/src/app/image-api-server /
COPY --from=builder-web /web/dist/ /web/dist/

USER nobody
EXPOSE 8080
ENV IMAGES_PATH /img

CMD ["/image-api-server"]
