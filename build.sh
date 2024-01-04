#!/bin/zsh
export GOOS=linux
export GOARCH=amd64
go build -o lark main.go
tag=${tag:-"latest"}
echo 'Docker Tag = '$tag
docker build -f dockerfile -t ydq1234/drone-lark:$tag .
docker push ydq1234/drone-lark:$tag
docker image prune -f