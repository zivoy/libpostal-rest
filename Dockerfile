FROM sleekybadger/libpostal:1.1-alpha-alpine AS libpostal-build
FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY --from=libpostal-build /usr/lib/libpostal.so /usr/lib/libpostal.so
COPY --from=libpostal-build /usr/lib/libpostal.so.1 /usr/lib/libpostal.so.1
COPY --from=libpostal-build /usr/include/libpostal /usr/include/libpostal
COPY --from=libpostal-build /usr/lib/pkgconfig/libpostal.pc /usr/lib/pkgconfig/libpostal.pc

RUN ldconfig /etc/ld.so.conf.d
RUN apk add --no-cache alpine-sdk

COPY go.sum go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -ldflags='-s -w' -trimpath -o /dist/app

RUN ldd /dist/app | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /dist/%
RUN ln -s ld-musl-x86_64.so.1 /dist/lib/libc.musl-x86_64.so.1

FROM scratch

COPY --from=libpostal-build /data /data

COPY --from=builder /dist/ /

EXPOSE 8724
USER 65534
CMD [ "/app" ]
