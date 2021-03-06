# gumdrop

## About

`gumdrop` is meant to provide a simple HTTP interface for dropping files onto a machine.  It's basically an FTP `PUT` operation when you don't have or want FTP.

Its API is simple:

* A single `POST` endpoint at `/`, that returns `201` when successful
  * Protected by an `Authorization: bearer Token` header (see parameters)
  * Expects a `Content-Type: multipart/form-data` header
    * With a `file` property that contains the file to be uploaded
  * Expects a `x-directory: someDirectory` header
    * Such that `file` will be stored at `x-directory/file`

## A Quick Run

```shell script
$ git clone https://github.com/charlesread/gumdrop.git
$ cd gumdrop
$ go run gumdrop.go
2020/10/16 07:33:43 Starting `gumdrop`...
2020/10/16 07:33:43 Address: ":8080"
2020/10/16 07:33:43 BaseDir: .
2020/10/16 07:33:43 FileMode: 644
2020/10/16 07:33:43 LogFilePath: 
2020/10/16 07:33:43 PathMode: 755
2020/10/16 07:33:43 Tokens: [superSecretToken someOtherEquallySuperSecretToken]
...
```

## Building

NOTE: This is a `go` application.  Please have `go` >= 1.15 installed before building.

```shell script
$ git clone https://github.com/charlesread/gumdrop.git
$ cd gumdrop
$ sudo make build
```

This will build `bin/gumdrop`.

## Installing

To build and install `gumdrop` at `/usr/local/bin/gumdrop` simply:

```shell script
$ sudo make install
``` 

### Install Service

`gumdrop` runs just fine from the command line, a `systemd` unit file is provided at `gumdrop.service`.

The default service runs as the user `gumdrop`, be sure that user exists:

```shell script
$ sudo useradd gumdrop -s /sbin/nologin -m
$ sudo cp config.yaml /home/gumdrop # copy over default config for convenience (optional)
$ sudo chown gumdrop:gumdrop /home/gumdrop/config.yaml
```

You can install the service with the `service` target:

```shell script
$ sudo make install
$ sudo make service
$ sudo journalctl -u gumdrop
```

### Docker

This repository includes a basic `Dockerfile` with `gumdrop` injected:

```shell script
$ docker build . -t gumdrop && \
  docker run -it -p 8080:8080 --name gumdrop --rm gumdrop
...  
Successfully built 3bdc5897174c
Successfully tagged gumdrop:latest
2020/10/16 17:56:44 Starting `gumdrop`...
2020/10/16 17:56:44 Address: ":8080"
2020/10/16 17:56:44 BaseDir: .
2020/10/16 17:56:44 FileMode: 644
2020/10/16 17:56:44 LogFilePath: 
2020/10/16 17:56:44 PathMode: 755
2020/10/16 17:56:44 Tokens: [superSecretToken someOtherEquallySuperSecretToken]
```

Note that the service will run as the `gumdrop` user and has the working directory as `/home/gumdrop`, thus keeping the default value of `BaseDir` (`.`) will drop files to `/home/gumdrop`.

## Configuration

Runtime configuration is done via a YAML file. This file must be named `config.yaml`, a sample is in the repository root.

This file can be placed in the following locations:

* `.`
* `/etc/gumdrop/`
* `$HOME/.gumdrop/`
* `$HOME`

### Passing Configuration Parameters

Parameters may also be set/overridden via environment variables, `GUMDROP_ADDRESS` for example.

### Parameters

| Name | Type | Default Value | Function | Environment Variable Override |
| ---- | ---- |-------------- | -------- | ----------------------------- |
| `Address` | string | `:8080` | Sets the address where `gumdrop` will serve. | `GUMDROP_ADDRESS` |
| `BaseDir` | string | `.` | The base directory where files will be dropped. | `GUMDROP_BASEDIR` |
| `FileMode` | uint32 | 0644 (in octal) | The file permissions, before umask, of the created file. | `GUMDROP_FILEMODE` |
| `LogFilePath` | string | `""` | The location of the log file where you'd like to log. By default, `gumdrop` will log to `os.Stdout` (STDOUT). | `GUMDROP_LOGFILEPATH` |
| `PathMode` | uint32 | 0755 (in octal) | The directory permissions, before umask, of the created file. | `GUMDROP_PATHMODE` |
| `Tokens` | []string | `[superSecretToken someOtherEquallySuperSecretToken]` | Tokens allowed in `Authorization: bearer Token` header. | _not available_ | 


## Running

`gumdrop` is entirely self-contained, simply run the executable:

```shell script
$ ./gumdrop
2020/10/16 07:33:43 Starting `gumdrop`...
2020/10/16 07:33:43 Address: ":8080"
2020/10/16 07:33:43 BaseDir: .
2020/10/16 07:33:43 FileMode: 644
2020/10/16 07:33:43 LogFilePath: 
2020/10/16 07:33:43 PathMode: 755
2020/10/16 07:33:43 Tokens: [superSecretToken someOtherEquallySuperSecretToken]
...
```

## Dropping Files

A single file:

```shell script
$ echo "some text" > someFile.txt
$ curl -X POST \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: bearer superSecretToken" \
  -H "x-directory: someDirectory" \
  -F file=@someFile.txt \
  localhost:8080
$ rm someFile.txt
```

A handful of files:

```shell script
$ curl -X POST \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: bearer superSecretToken" \
  -H "x-directory: someDirectory" \
  -F file=@someMassiveFile0.zip \
  -F file=@someMassiveFile1.zip \
  -F file=@someMassiveFile2.zip \
  -F file=@someMassiveFile3.zip \
  localhost:8080
```

## Removal

```shell script
$ sudo make remove
```

Or replicate the commands in the `remove` Makefile target.

## The TL;DR Installation

If you're on a `systemd` system, have:

* `go`
* `git`
* `useradd`
* `make`

Then:

```shell script
$ cd /tmp
$ git clone https://github.com/charlesread/gumdrop.git
$ cd gumdrop            
$ sudo useradd gumdrop -s /sbin/nologin -m
$ sudo cp config.yaml /home/gumdrop
$ sudo chown gumdrop:gumdrop /home/gumdrop/config.yaml # edit appropriately
$ sudo make install
$ sudo make service
```

OR

A simple shell script exists that will make the user, copy the default config to `/home/gumdrop`, install `gumdrop` and its `systemd` service, and start it (literally exactly what is above).

```shell script
$ cd /tmp
$ git clone https://github.com/charlesread/gumdrop.git
$ cd gumdrop   
$ ./install.sh 
```
