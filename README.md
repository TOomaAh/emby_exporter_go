# Emby Exporter

An exporter for emby that exports Emby's main metrics.


## For build this project. It's simple
`go get &&
go build .`


## To build the docker image (it's already build but just if you want):

To build the exporter you need a maxmind database. To do this, go to [Maxmind](https://www.maxmind.com/en/home) and download the GeoLite2 City database.

The ".mmdb" file should then be placed at the root of the project with the following name: `geoip.mmdb`


`docker build -t emby_exporter .`

### RUN Docker container

```
docker run -d -it \
   --name=emby_exporter \
   -e TZ=Europe/Paris \
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
options: # optional
  geoip: true # optional : default false
```

### How to get your userID

1. Login to your emby server
2. Go to the settings page
3. Go to the profile page
4. In the url, you will find your userID
    (ex: http://localhost:8096/web/index.html#!/settings/profile.html?userId=<YOUR_USER_ID>&serverId=xxxxx)
    
### Get metrics

You can access the metrics on the following url:
`http://ip:9210/metrics`

## Grafana Dashboard
(I will share my dashboard later)

![Dashboard example](https://github.com/TOomaAh/emby_exporter_go/blob/main/example/dashboard_grafana.png)

This project will end up in a docker to facilitate its use.
