#!/bin/sh

cd path/to/godiscourse/web || exit
cp .env.production .env || exit
rm -rf dist || exit
npm run build || exit
rsync -rcv dist/* remote.server.host:path/to/godiscourse/html/ || exit
cp .env.development .env || exit
