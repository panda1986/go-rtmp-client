#!/bin/bash
echo 'start build'
go build -o rtmp_client ./src/ || exit 1
echo 'build finish, run ./rtmp_client'