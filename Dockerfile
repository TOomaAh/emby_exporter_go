FROM golang:1.16-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o emby_exporter
RUN apk update && apk add dos2unix && dos2unix entrypoint.sh


FROM alpine:latest

WORKDIR /app

COPY --from=0 /app/emby_exporter /app/emby_exporter
COPY --from=0 /app/entrypoint.sh /app/entrypoint.sh
COPY --from=0 /app/geoip.mmdb /app/geoip.mmdb

RUN chmod +x entrypoint.sh
RUN mkdir /config && touch /config/config.yml

EXPOSE 9210

ENTRYPOINT [ "./entrypoint.sh" ]