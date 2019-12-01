BUILD_DIR = build
SERVICE = agent
CGO_ENABLED ?= 0
GOOS ?= linux

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o ${BUILD_DIR}/mainflux-$(1) cmd/main.go
endef

all: $(SERVICE)

.PHONY: all $(SERVICE) dockers dockers_dev

clean:
	rm -rf ${BUILD_DIR}

install:
	cp ${BUILD_DIR}/* $(GOBIN)

test:
	go test -v -race -count 1 -tags test $(shell go list ./... | grep -v 'vendor\|cmd')

$(SERVICE):
	$(call compile_service,$(@))

run:
	cd $(BUILD_DIR) && ./mainflux-$(SERVICE)
