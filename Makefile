PREFIX="$(HOME)"
BIN="$(PREFIX)/bin"

all:
	go build -o build/rootless-orchestration main.go
clean:
	rm -rf build/
install:
	mkdir -p "$(BIN)"
	cp build/rootless-orchestration "$(BIN)/rootless-orchestration"
uninstall:
	rm -rf "$(BIN)/rootless-orchestration"
