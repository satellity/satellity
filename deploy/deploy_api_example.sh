#!/bin/sh

cd path/to/godiscourse || exit
sed -i ''  "s/BUILD_VERSION/`git rev-parse HEAD`/g" internal/configs/config.go || exit
make production || exit
ssh remote.server.host mv /path/to/godiscourse /path/to/godiscourse.old || exit
scp path/to/godiscourse/bin/godiscourse remote.server.host:godiscourse/godiscourse || exit
