#!/bin/sh

if [ -z "${CONFIG_FILE}" ]; then
    exec ./emby_exporter
else
    exec ./emby_exporter -c $CONFIG_FILE
fi