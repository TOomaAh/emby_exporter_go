#!/bin/sh
if [ -z "${CONFIG_FILE}" ]; then
    ./emby_exporter
else
    ./emby_exporter -c $CONFIG_FILE
fi
