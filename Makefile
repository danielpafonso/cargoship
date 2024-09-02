.SILENT:

.PHONY: build clean shipper loader packager

FLAGS = -trimpath -a -ldflags '-w -s'

all: clean build

define build_and_copy
	go build $(FLAGS) -o ./build/$(1) ./cmd/$(1)

	echo "$(2)Copy configuration file"
	cp configs/$(1)_config.yaml build/
endef

build:
	mkdir -p build

	echo "Build application"
	echo "- Shipper"
	$(call build_and_copy,shipper," ")
	echo "- Loader"
	$(call build_and_copy,loader," ")
	echo "- Packager"
	$(call build_and_copy,packager," ")

shipper:
	mkdir -p build

	echo "Build Shipper application "
	$(call build_and_copy,shipper)

loader:
	mkdir -p build

	echo "Build Loader application "
	$(call build_and_copy,loader)

packager:
	mkdir -p build

	echo "Build Packager application "
	$(call build_and_copy,packager)

clean:
	echo "Cleaning artifact folder"
	if [ -d build ]; then rm -rf build/*; fi
