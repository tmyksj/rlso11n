PREFIX=$(HOME)

_BIN=$(PREFIX)/bin
_URI_DOCKER=https://download.docker.com/linux/static/test/$(shell uname -m)/docker-20.10.0-rc1.tgz

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
		git clone https://github.com/rootless-containers/rootlesskit.git -b v0.10.1; \
		cd rootlesskit; \
		make; \
		cp bin/* $(_BIN); \
		\
		cd $$temp; \
		git clone https://github.com/rootless-containers/slirp4netns.git -b v0.4.7; \
		cd slirp4netns; \
		./autogen.sh; \
		./configure --prefix=$(PREFIX); \
		make; \
		make install; \
		\
		cd $$temp; \
		if which nvidia-container-runtime-hook > /dev/null 2>&1 && \
				[ -f /etc/nvidia-container-runtime/config.toml ] && \
				[ ! -f $(_BIN)/nvidia-container-runtime-hook ]; then \
			echo "#!/bin/sh" >> $(_BIN)/nvidia-container-runtime-hook; \
			echo "$$(which nvidia-container-runtime-hook) -config=$(_BIN)/nvidia-container-runtime-hook.config \"\$$@\"" >> $(_BIN)/nvidia-container-runtime-hook; \
			chmod +x $(_BIN)/nvidia-container-runtime-hook; \
			sed -e "s/#no-cgroups = false/no-cgroups = true/" /etc/nvidia-container-runtime/config.toml >> $(_BIN)/nvidia-container-runtime-hook.config; \
		fi; \
	'
uninstall:
	rm -rf $(_BIN)/rlso11n

_mkdir_bin:
	mkdir -p $(_BIN)
