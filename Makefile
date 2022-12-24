Version := $(shell git describe --tags --dirty || echo "v0.0.0")
LDFLAGS := "-w -s -X main.Version=$(Version)"

build:
	go build -ldflags $(LDFLAGS) -o bin/kube-create-secret ./cmd/main.go

install: build
	mv bin/kube-create-secret ~/bin/