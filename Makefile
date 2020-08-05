OSES    := linux
ARCHS   := amd64
SRCS    := $(wildcard *.go)
VER     := $(shell grep -Eo 'VERSION = `(.*)`' main.go | cut -d'`' -f2)
TGTS    := $(foreach os,$(OSES),$(foreach arch,$(ARCHS),bin/imghitcntr-$(os)-$(arch)))
BUILD   := $(shell echo `whoami`@`hostname -s` on `date`)
LDFLAGS := -ldflags='-X "main.build=$(BUILD)"'

.PHONY: clean dev test benchmark

all: $(TGTS) bin/checksums.md5

test: $(TGTS)

clean:
	@rm -f bin/*

benchmark:
	go test -bench . -benchmem

dev: $(SRCS)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o imghitcntr-$@-$(VER) .

$(TGTS): $(SRCS)
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$@)) go build $(LDFLAGS) -o $@-$(VER) .

$(SRCS):

bin/checksums.md5:
	cd bin && md5sum * > checksums.md5
