#!/bin/bash

VERSION="0.1.4"
URL="https://github.com/sosedoff/lunchy-go/releases/download/v${VERSION}/lunchy"
BIN_PATH="/usr/local/bin/lunchy"

echo "Downloading and installing lunchy v${VERSION}"
wget -q -O $BIN_PATH $URL && chmod +x $BIN_PATH
echo "Done. Installed into ${BIN_PATH}"