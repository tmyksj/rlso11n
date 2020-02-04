PREFIX=$(HOME)

_BIN=$(PREFIX)/bin
_URI_DOCKER=https://master.dockerproject.org/linux/x86_64/docker.tgz

all:
	go build -o build/rlso11n main.go
clean:
	rm -rf build/
install: _mkdir_bin
	cp build/rlso11n $(_BIN)/rlso11n
install_dependencies: _mkdir_bin
	/bin/sh -c '\
		set -e -x; \
		\
		temp=$$(pwd)/build/temp; \
		mkdir -p $$temp; \
		trap "rm -rf $$temp" EXIT INT TERM; \
		\
		cd $$temp; \
		curl -L -o docker.tgz $(_URI_DOCKER); \
		tar zxf docker.tgz --directory=$(_BIN) --strip-components=1; \
		\
		cd $$temp; \
		git clone https://github.com/rootless-containers/rootlesskit.git -b v0.7.0; \
		cd rootlesskit; \
		make; \
		cp bin/* $(_BIN); \
		\
		cd $$temp; \
		git clone https://github.com/rootless-containers/slirp4netns.git -b v0.4.2; \
		cd slirp4netns; \
		./autogen.sh; \
		./configure --prefix=$(PREFIX); \
		make; \
		make install; \
	'
uninstall:
	rm -rf $(_BIN)/rlso11n

_mkdir_bin:
	mkdir -p $(_BIN)
