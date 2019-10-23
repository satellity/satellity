#!/bin/sh

cd path/to/satellity/app || exit
yarn install || exit
rm -rf build || exit
yarn build || exit
rsync -rcv --delete build/ remote.server.host:path/to/satellity/html/ || exit
