FROM golang:1.26-alpine

WORKDIR /app
COPY . .
RUN go mod download
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o emby_exporter ./cmd/emby_exporter_go
RUN apk update && apk add dos2unix tzdata && dos2unix entrypoint.sh


FROM alpine:latest

WORKDIR /app

COPY --from=0 /app/emby_exporter /app/emby_exporter
COPY --from=0 /app/entrypoint.sh /app/entrypoint.sh
COPY --from=0 /app/geoip.mmdb /app/geoip.mmdb

RUN apk update && apk add tzdata && rm -rf /var/cache/apk/*

RUN chmod +x entrypoint.sh
RUN mkdir /config && touch /config/config.yml

EXPOSE 9210

ENTRYPOINT [ "./entrypoint.sh" ]
