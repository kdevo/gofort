GO ?= go
SHELL = /bin/sh

# Build
BUILD_DIR ?= ./bin
EXECUTABLE := gofort
VERSION ?= $(shell git describe --tags --abbrev=0)
VERSION := $(or ${VERSION},v0.1.0)
GOOSS = linux darwin windows
GOARCHS = 386 amd64 arm arm64
COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS = "-s -w \
	-X 'main.NAME=${EXECUTABLE}' \
	-X 'main.VERSION=${VERSION}' \
	-X 'main.COMMIT=${COMMIT}' \
	-X 'main.OS={OS}' \
	-X 'main.ARCH={ARCH}' \
	-X 'main.DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")'"

# Fortune
FORTUNE_SRC ?= $(if $(wildcard /usr/share/fortune/.),/usr/share/fortune,/usr/share/games/fortunes)
FORTUNE_DST = ./pkg/fortune/texts
FORTUNE_TXT := $(patsubst %.dat,%,$(wildcard ${FORTUNE_SRC}/*.dat))

.PHONY: fortune test build release clean
all: fortune test build release

fortune: ${FORTUNE_TXT}
	@mkdir -p ${FORTUNE_DST}
	@cp  ${FORTUNE_TXT} ${FORTUNE_DST}

test: fortune
	${GO} test -v -race ./... 

build: fortune
	${GO} build -o ${BUILD_DIR}/${EXECUTABLE} \
		-ldflags=$(subst {ARCH},$(shell go env GOARCH),$(subst {OS},$(shell go env GOOS),${LDFLAGS}) ./cmd/gofort) 

release: fortune
# Crossing architectures and OSs for cross-compiling:
	$(foreach arch,${GOARCHS}, \
		$(foreach os,${GOOSS}, \
			$(shell GOARCH=${arch} GOOS=${os} \
				${GO} build -o ${BUILD_DIR}/${EXECUTABLE}-${VERSION}-${os}_${arch} \
					-ldflags=$(subst {ARCH},${arch},$(subst {OS},${os},${LDFLAGS}) ./cmd/gofort)) \
		) \
	) 
# Add exe to windows binaries for easier execution:
	$(foreach winbin,$(filter-out $(wildcard ${BASENAME}*.exe),$(wildcard ${BASENAME}-windows*)), \
		$(shell mv ${winbin} ${winbin}.exe) \
	)

clean:
	rm bin/${EXECUTABLE}*
