#!/bin/bash
set -e

function getversion() {
    branch=`git rev-parse --abbrev-ref HEAD`
    commitid=`git rev-parse --short HEAD`
    builddate=`date +%Y%m%d-%H%M%S`
    echo $branch-$commitid-$builddate
}

cd `dirname $0`

# use go vendor
export GO15VENDOREXPERIMENT=1

version=`getversion`
pkgpath="github.com/4paradigm/cfg-center/src/cfg-server"

echo "build cfg-server-linux"
GOOS=linux GOARCH=amd64  go build -ldflags "-X main.versionStr=${version}"  -o bin/cfg-server-linux ${pkgpath}

#echo "build cfg-server-osx"
#GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.versionStr=${version}"  -o bin/cfg-server-osx ${pkgpath}

#echo "build cfg-server-windows"
#GOOS=windows GOARCH=386  go build -ldflags "-X main.versionStr=${version}"  -o bin/cfg-server-win ${pkgpath}

echo "done"

