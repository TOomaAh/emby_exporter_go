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
   -e CONFIG_FILE=NAME_OF_YOUR_FILE.yml \
   -v '/path/to/your/config/file.yml:/config/file.yml' \
   bagul/goemby_exporter:latest
```

### Config file example
```yaml
server:
  url: "http://<ip|domain name>"
  port: 8096
  token: "your token"
  userID: "your userID"
```

## Grafana Dashboard
(I will share my dashboard later)

![Dashboard example](https://github.com/TOomaAh/emby_exporter_go/blob/main/example/dashboard_grafana.png)

This project will end up in a docker to facilitate its use.
