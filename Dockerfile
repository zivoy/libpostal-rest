FROM sleekybadger/libpostal:1.1-alpha-alpine AS libpostal-build
FROM golang:1.23-alpine AS healthcheck-build
WORKDIR /build

# Build a small static binary for healthcheck
RUN printf 'package main\nimport (\n\t"net/http"\n\t"os"\n)\nfunc main() {\n\tresp, err := http.Get("http://localhost:8724/")\n\tif err != nil || resp.StatusCode != 200 { os.Exit(1) }\n\tos.Exit(0)\n}' > healthcheck.go && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o healthcheck healthcheck.go

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

ARG VERSION=latest

COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w -X main.Version=$VERSION" -trimpath -o /dist/app

RUN ldd /dist/app | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /dist/%
RUN ln -s ld-musl-x86_64.so.1 /dist/lib/libc.musl-x86_64.so.1
COPY --from=healthcheck-build /build/healthcheck /dist/healthcheck

FROM scratch

COPY --from=libpostal-build /data /data

COPY --from=builder /dist/ /

EXPOSE 8724
USER 65534

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/healthcheck"]

CMD [ "/app" ]
