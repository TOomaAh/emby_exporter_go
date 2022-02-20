FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o emby_exporter
RUN chmod +x entrypoint.sh

EXPOSE 9210

ENTRYPOINT [ "./entrypoint.sh" ]