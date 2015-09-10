#!/bin/bash

VERSION="0.2.0"
URL="https://github.com/sosedoff/lunchy-go/releases/download/v${VERSION}/lunchy"
BIN_PATH="/usr/local/bin/lunchy"

if [ -e $BIN_PATH ]
then
  echo "Removing already installed version"
  rm $BIN_PATH
fi

echo "Downloading and installing lunchy v${VERSION}"
curl -sL $URL -o $BIN_PATH && chmod +x $BIN_PATH
echo "Done. Installed into ${BIN_PATH}"