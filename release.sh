#!/bin/sh

set -e

FILEPREFIX="plaintodo"
BUILD_COMMAND="gom build"
BUILD_DIR="plaintodo/*.go"
BUILD_OPTION=""

OS=("darwin" "windows" "linux")
ARCH=("amd64")

for GOOS in ${OS[@]};
do
  for GOARCH in ${ARCH[@]};
  do
    FILENAME=${FILEPREFIX}_${GOOS}_${GOARCH}
    if test $GOOS = "windows" ; then
      FILENAME=${FILENAME}.exe
    fi

    echo "GOOS=${GOOS} GOARCH=${GOARCH} FILENAME=${FILENAME}"
    GOOS=${GOOS} GOARCH=${GOARCH} $BUILD_COMMAND -o $FILENAME $BUILD_DIR $BUILD_OPTION
    zip ${FILENAME}.zip ${FILENAME}
    rm ${FILENAME}
  done
done
