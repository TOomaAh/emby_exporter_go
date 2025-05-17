# Emby Exporter

An exporter for emby that exports Emby's main metrics.

## /!\ IMPORTANT /!\

This project may no longer be maintained. I no longer use emby because of the 25-device limit, which I have unfortunately exceeded. 

As I couldn't find a solution with the Emby team and couldn't afford to take out a subscription just to extend the limit, I was forced to stop using emby.

If you ever want to add features, I invite you to create a fork of the project and either continue it on your own, or make a PR that I'll validate if I can.

I will continue to respond to issues and try to solve your problems or improve it as much as possible.

Thank you for your understanding.

### Get metrics

You can access the metrics on the following url:
`http://ip:9210/metrics`

## Grafana Dashboard

![Dashboard example](https://github.com/TOomaAh/emby_exporter_go/blob/main/example/dashboard_grafana.png)

[Dashboard link](https://github.com/TOomaAh/emby_exporter_go/blob/main/example/Emby.Dashboard-1703419734858.json)

I thank [jaycedk](https://github.com/jaycedk) for the dashboard (it's his)


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
  geoipOptions: # optional if geoip is true
    accountId: "your maxmind account id" # optional
    licenseKey: "your maxmind license key" # optional
```

## MaxMind GeoIP Database

The exporter supports GeoIP lookups to provide geographic information about client connections. This requires a MaxMind GeoLite2 City database.

For more detailed information on MaxMind database updates, see their [developer documentation](https://dev.maxmind.com/geoip/updating-databases/#directly-downloading-databases).

### Options for using GeoIP:

1. **Manual database download**:
   - Register for a free account at [MaxMind](https://www.maxmind.com/en/geolite2/signup)
   - Download the GeoLite2 City database (.mmdb file)
   - Place it at the root of the project or in your Docker volume as `geoip.mmdb`
   - Enable GeoIP in your config with `geoip: true`

2. **Automatic database updates**:
   - Register for a free account at [MaxMind](https://www.maxmind.com/en/geolite2/signup)
   - Find your Account ID ([How to find your Account ID](https://support.maxmind.com/hc/en-us/articles/4412951066779-Find-my-Account-ID))
   - Generate or access your License Key ([MaxMind License Keys page](https://www.maxmind.com/en/accounts/current/license-key))
   - Add these credentials to your config file:
   ```yaml
   options:
     geoip: true
     geoipOptions:
       accountId: "your_account_id_here"
       licenseKey: "your_license_key_here"
   ```
   - The exporter will automatically download and update the GeoIP database

### Environment Variables

You can also specify a custom path for your GeoIP database using an environment variable:

```
GEOIP_DB=/path/to/your/GeoLite2-City.mmdb
```

In a Docker container:
```
docker run -d -it \
   --name=emby_exporter \
   -e TZ=Europe/Paris \
   -e CONFIG_FILE=NAME_OF_YOUR_FILE.yml \
   -e GEOIP_DB=/config/GeoLite2-City.mmdb \
   -v '/path/to/your/config/:/config/' \
   bagul/goemby_exporter:latest
```

### How to get your userID

1. Login to your emby server
2. Go to the settings page
3. Go to the profile page
4. In the url, you will find your userID
    (ex: http://localhost:8096/web/index.html#!/settings/profile.html?userId=<YOUR_USER_ID>&serverId=xxxxx)