OS  := $(shell uname -s)
ENV := CGO_LDFLAGS_ALLOW="-s|-w"

ifeq ($(OS),Darwin)
ENV += PKG_CONFIG_PATH="$(shell brew --prefix libffi)/lib/pkgconfig" CGO_CFLAGS_ALLOW="-Xpreprocessor"
endif

.PHONY: build test
default: build

build:
	$(ENV) go build .

test:
	$(ENV) go test .

converter:
	$(ENV) go build -o bin/converter examples/convert/convert.go
crop:
	$(ENV) go build -o bin/crop examples/crop/crop.go
watermark:
	$(ENV) go build -o bin/watermark examples/watermark/watermark.go