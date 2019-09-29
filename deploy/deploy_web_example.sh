#!/bin/sh

cd path/to/satellity/web || exit
cp .env.production .env || exit
rm -rf dist || exit
npm run build || exit
rsync -rcv --delete dist/ remote.server.host:path/to/satellity/html/ || exit
cp .env.development .env || exit
