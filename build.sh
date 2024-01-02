#!/bin/zsh
export GOOS=linux
export GOARCH=amd64
tag=${tag:-"dev"}
echo $tag
go build -o lark main.go
docker build -f dockerfile -t ydq1234/drone-lark:$tag .
docker push ydq1234/drone-lark:$tag
docker image prune -f