FROM golang:1.14.2-buster AS builder

RUN apt-get update -y && apt-get install -y libpng-dev

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM gcr.io/distroless/cc-debian10

COPY --from=builder /go/bin/alpacascloud /
COPY --from=builder /usr/lib/x86_64-linux-gnu/libpng16.so.16 /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libpng16.so.16.36.0 /usr/lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libz.so.1 /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libz.so.1.2.11 /lib/x86_64-linux-gnu/

USER nobody
EXPOSE 8080
ENV IMAGES_PATH /img

CMD ["/alpacascloud"]
