#!/bin/sh

cd path/to/satellity || exit
sed -i ''  "s/BUILD_VERSION/`git rev-parse HEAD`/g" internal/configs/config.go || exit
make production || exit
ssh remote.server.host mv /path/to/satellity /path/to/satellity.old || exit
scp path/to/satellity/bin/satellity remote.server.host:satellity/satellity || exit
