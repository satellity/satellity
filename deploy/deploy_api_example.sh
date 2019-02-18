#!/bin/sh

cd $GOPATH/src/github.com/godiscourse/godiscourse/api
cp config/production.cfg config/config.go
sed -i ''  "s/BUILD_VERSION/`git rev-parse HEAD`/g" config/config.go || exit
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
cp config/dev.cfg config/config.go
ssh godiscourse mv /path/to/godiscourse /path/to/godiscourse.old
scp api godiscourse:path/to/godiscourse
