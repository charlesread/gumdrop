all: build

build: clean
	@echo "Building bin/linux and bin/windows"
	@GOOS=linux GOARCH=386 go build -o bin/linux/gumdrop gumdrop.go
	@GOOS=windows GOARCH=386 go build -o bin/windows/gumdrop.exe gumdrop.go

install: remove clean build
	@echo "Moving binary to /usr/local/bin/gumdrop"
	@mkdir -p /usr/local/bin
	@cp bin/linux/gumdrop /usr/local/bin/gumdrop
	@chmod 744 /usr/local/bin/gumdrop

service:
	@echo "Coming soon..."

clean:
	@echo "Removing bin/**"
	@rm -rf bin/**

remove:
	@echo "Removing /usr/local/bin/gumdrop"
	@rm -rf /usr/local/bin/gumdrop