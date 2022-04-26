OS  := $(shell uname -s)
ENV := CGO_LDFLAGS_ALLOW="-s|-w"

ifeq ($(OS),Darwin)
ENV += PKG_CONFIG_PATH="$(shell brew --prefix libffi)/lib/pkgconfig" CGO_CFLAGS_ALLOW="-Xpreprocessor"
endif

.PHONY: all
all: test converter crop empty watermark

.PHONY: test
test:
	$(ENV) go test .

.PHONY: converter
converter:
	$(ENV) go build -o bin/converter examples/convert/convert.go

.PHONY: crop
crop:
	$(ENV) go build -o bin/crop examples/crop/crop.go

.PHONY: watermark
watermark:
	$(ENV) go build -o bin/watermark examples/watermark/watermark.go

.PHONY: empty
empty:
	$(ENV) go run examples/empty_img/empty.go

.PHONY: to-pdf
to-pdf:
	$(ENV) go run examples/to_pdf/to_pdf.go
