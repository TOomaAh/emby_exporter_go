# Emby Exporter

Hi there!
Normally the project works. I invite you to try it, I have to refine it a bit to get some logs of things like that. I invite you to try it and if there is a problem or if you have a request don't hesitate to write an issue :)


## For build this project. It's simple
`go get &&
go build .`


## To build the docker image:

`docker build -t emby_exporter .`
```
docker run -d -it \
   --name=emby_exporter \
   -p 9210:9210 \
   -e PORT=9210 \ #OPTIONAL
   -e SCHEME=http \
   -e EMBYURL=localhost \
   -e EMBYPORT=8096 \
   -e USERID=youruserid \
   -e TOKEN=yourembytoken \
   bagul/goemby_exporter:latest
```

## Grafana Dashboard
(I will share my dashboard later)

![Dashboard example](https://github.com/TOomaAh/emby_exporter_go/blob/main/example/dashboard_grafana.png)

This project will end up in a docker to facilitate its use.
