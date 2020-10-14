all: build

build: clean
	@echo "Building bin/gumdrop"
	@go build -o bin/gumdrop gumdrop.go

install: remove clean build
	@echo "Moving binary to /usr/local/bin/gumdrop"
	@mkdir -p /usr/local/bin
	@cp bin/gumdrop /usr/local/bin/gumdrop
	@chmod 755 /usr/local/bin/gumdrop

service:
	@echo "Coming soon..."

clean:
	@echo "Removing bin/**"
	@rm -rf bin/**

remove:
	@echo "Removing /usr/local/bin/gumdrop"
	@rm -rf /usr/local/bin/gumdrop