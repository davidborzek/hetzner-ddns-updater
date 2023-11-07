FROM golang:1.21.4-alpine3.17 AS base

RUN adduser -D -H hetzner-ddns-updater

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go mod download

FROM base as build

RUN go build -o hetzner-ddns-updater -tags prod main.go

FROM scratch as prod

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /build/hetzner-ddns-updater /

USER hetzner-ddns-updater:hetzner-ddns-updater

CMD ["./hetzner-ddns-updater"]
