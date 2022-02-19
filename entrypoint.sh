#!/bin/bash
port=$PORT
if [ -z "$port" ]; then
    port=9210
fi
embyport=$EMBYPORT
if [ -z "$embyport" ]; then
    embyport=8096
fi

url=$EMBYURL
if [ -z "$URL" ]; then
    URL=http://localhost:$embyport
fi
embytoken=$TOKEN
if [ -z "$embytoken" ]; then
    echo "You must provide a token for the emby server"
    exit 1
fi

useridemby=$USERID
if [ -z "$useridemby" ]; then
    echo "You must provide a user id for the emby server"
    exit 1
fi

./emby-server --port $port --emby $url --embyport $embyport --token $embytoken --user-id $useridemby