FROM golang:1-alpine AS builder
RUN apk add --no-cache build-base

RUN apk update && \
   apk upgrade -U && \
   apk add libwebp-dev && \
   apk add --no-cache ca-certificates && \
   rm -rf /var/cache/* && \
   update-ca-certificates

WORKDIR /build
COPY go.work .
COPY services/api/go.mod ./services/api/
# COPY services/api/go.sum ./services/api/
COPY internal/db/go.mod ./internal/db/
COPY internal/db/go.sum ./internal/db/
COPY services/tracker/go.mod ./services/tracker/
COPY services/tracker/go.sum ./services/tracker/
RUN go work sync

COPY . .
WORKDIR /build/services/api
RUN go build -ldflags='-s -w' -trimpath -o /dist/app main.go
RUN ldd /dist/app | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /dist/%
RUN ln -s ld-musl-x86_64.so.1 /dist/lib/libc.musl-x86_64.so.1

FROM scratch
COPY --from=builder /dist /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER 65534
EXPOSE 8080
WORKDIR /
ENTRYPOINT ["/app"]