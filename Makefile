all: build

build: test clean
	@echo "\n*** Building bin/gumdrop ***\n"
	@go build -o bin/gumdrop gumdrop.go

install: remove clean build
	@echo "\n*** Moving binary to /usr/local/bin/gumdrop ***\n"
	@mkdir -p /usr/local/bin
	@cp bin/gumdrop /usr/local/bin/gumdrop
	@chmod 755 /usr/local/bin/gumdrop

dist:
	-rm -r dist
	GOOS=linux GOARCH=amd64 go build -o dist/linux_amd64/gumgrop gumdrop.go
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin_amd64/gumgrop gumdrop.go
	GOOS=windows GOARCH=amd64 go build -o dist/windows_amd64/gumgrop.exe gumdrop.go
	cp config.yaml dist/linux_amd64
	cp config.yaml dist/darwin_amd64
	cp config.yaml dist/windows_amd64
ifdef v
	cd dist && tar -czvf gumdrop-$(v).tar.gz --exclude=*.tar* .
else
	cd dist && tar -czvf gumdrop.tar.gz --exclude=*.tar* .
endif

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

integration:
	@echo "\n*** Beginning integration testing ***\n"
	@echo "\n*** Building containter ***\n"
	-docker build . -q -t gumdrop_integration
	@echo "\n*** Running container ***\n"
	-docker run -d -p 8888:8080 --name gumdrop_integration --rm gumdrop_integration
	# let's run the integration test _in_ the container
	-docker exec -it gumdrop_integration /usr/local/go/bin/go test -tags integration -v internal/integration_test/integration_test.go
	@echo "\n*** Stopping container ***\n"
	-docker stop gumdrop_integration

test:
	@echo "\n*** Running tests ***\n"
	@go test -v ./...

test_all: test integration