#!/bin/sh
port=$PORT
if [ -z "$port" ]; then
    port=9210
fi
embyport=$EMBYPORT
if [ -z "$embyport" ]; then
    embyport=8096
fi

scheme=$SCHEME
if [ -z "$scheme" ]; then
    scheme=http
fi

if [ -z "$EMBYURL" ]; then
    url="${scheme}://localhost"
else
    url=${scheme}://${EMBYURL}
fi
embytoken=$TOKEN
if [ -z "${embytoken}" ]; then
    echo "You must provide a token for the emby server"
    exit 1
fi

useridemby=$USERID
if [ -z "$useridemby" ]; then
    echo "You must provide a user id for the emby server"
    exit 1
fi
echo "Connecting to emby server at ${url}:${embyport}"
./emby_exporter --port $port --emby $url --embyport $embyport --token $embytoken --user-id $useridemby