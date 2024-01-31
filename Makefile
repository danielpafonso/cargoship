.SILENT:

.PHONY: build clean

all: clean build

build:
	mkdir -p build
	#
	echo "Build application"
	echo "- Loader"
	go build -trimpath -a -ldflags '-w -s' -o ./build/loader ./cmd/loader
	echo "- Shipper"
	go build -trimpath -a -ldflags '-w -s' -o ./build/shipper ./cmd/shipper
	#
	echo "Copy configuration file"
	cp configs/loader_config.yaml build/
	cp configs/shipper_config.yaml build/

clean:
	echo "Cleaning artifact folder"
	if [ -d build ]; then rm -rf build/*; fi
