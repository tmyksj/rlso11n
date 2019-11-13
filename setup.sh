#!/bin/sh

set -e -x

PREFIX="$HOME"
BIN="$PREFIX/bin"
URI_DOCKER="https://master.dockerproject.org/linux/x86_64/docker.tgz"
URI_ROOTLESS="https://master.dockerproject.org/linux/x86_64/docker-rootless-extras.tgz"

tmp=$(mktemp -d)
trap "rm -rf $tmp" EXIT INT TERM

mkdir -p "$BIN"

cd "$tmp"
curl -L -o docker.tgz "$URI_DOCKER"
cd "$BIN"
tar zxf "$tmp/docker.tgz" --strip-components=1

cd "$tmp"
curl -L -o rootless.tgz "$URI_ROOTLESS"
cd "$BIN"
tar zxf "$tmp/rootless.tgz" --strip-components=1

cd "$tmp"
mkdir -p "$tmp/go"
export GOPATH="$tmp/go"
go get github.com/rootless-containers/rootlesskit/cmd/rootlessctl
cd "$BIN"
cp "$tmp/go/bin/rootlessctl" rootlessctl

cd "$tmp"
git clone https://github.com/rootless-containers/slirp4netns.git -b v0.3.0
cd "$tmp/slirp4netns"
./autogen.sh
./configure --prefix="$PREFIX"
make
make install
