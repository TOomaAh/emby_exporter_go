FROM golang:1.26-alpine AS build

WORKDIR /src

# Download dependencies first so this layer is cached until go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
    -o /emby_exporter ./cmd/emby_exporter_go


FROM alpine:3.22

RUN apk add --no-cache tzdata && mkdir /config && touch /config/config.yml

WORKDIR /app

COPY --from=build /emby_exporter /app/emby_exporter
COPY geoip.mmdb /app/geoip.mmdb

EXPOSE 9210

ENTRYPOINT ["/app/emby_exporter"]
