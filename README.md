# Emby Exporter

Hi there!
Normally the project works. I invite you to try it, I have to refine it a bit to get some logs of things like that. I invite you to try it and if there is a problem or if you have a request don't hesitate to write an issue :)


## For build this project. It's simple
`go get &&
go build .`


## To build the docker image: 

`docker build -t emby_exporter .`
```
docker run \
   --name=emby_exporter \
   -p 9210:9210 \ #OPTIONAL
   -v PORT=9210 \
   -v EMBYURL=http://localhost \
   -v EMBYPORT=8096 \
   -v USERID=youruserid \
   -v TOKEN=yourembytoken \
   emby_exporter
```

This project will end up in a docker to facilitate its use.