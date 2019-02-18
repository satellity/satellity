#!/bin/sh

cd $GOPATH/src/github.com/godiscourse/godiscourse/web
rm -rf dist
npm run build
rsync -rcv dist/* godiscourse:path/to/godiscourse/html/
