all: build

test:
	@echo "\n*** Running tests ***\n"
	@go test -v ./...

build: test clean
	@echo "\n*** Building bin/gumdrop ***\n"
	@go build -o bin/gumdrop gumdrop.go

install: remove clean build
	@echo "\n*** Moving binary to /usr/local/bin/gumdrop ***\n"
	@mkdir -p /usr/local/bin
	@cp bin/gumdrop /usr/local/bin/gumdrop
	@chmod 755 /usr/local/bin/gumdrop

service:
	@echo "\n*** Copying unit file and enabling service ***\n"
	cp gumdrop.service /etc/systemd/system/
	systemctl enable gumdrop
	systemctl start gumdrop
	systemctl status gumdrop

clean:
	@echo "\n*** Removing bin/** ***\n"
	@rm -rf bin/**

remove:
	@echo "\n*** Removing components ***\n"
	-systemctl stop gumdrop
	-systemctl disable gumdrop
	-rm -rf /usr/local/bin/gumdrop
	-rm -f /etc/systemd/system/gumdrop.service
	-systemctl daemon-reload
	-systemctl reset-failed