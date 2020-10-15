# gumdrop

## About

`gumdrop` is meant to provide a simple HTTP interface for dropping files onto a machine.  It's basically an FTP `PUT` operation when you don't have or want FTP.

Its API is simple:

* A single `POST` endpoint at `/`, that returns `201` when successful
  * Protected by `Authorization: bearer Token` (see parameters)
  * Expects a `Content-Type: multipart/form-data`
    * With a `file` property that contains the file to be uploaded
  * Expects a `x-directory: someDirectory` header
    * Such that `file` will be stored at `x-directory/file`

## A Quick Run

```shell script
$ git clone https://github.com/charlesread/gumdrop.git
$ cd gumdrop
$ go run gumdrop.go
2020/10/14 19:44:36 Starting `gumdrop`...
2020/10/14 19:44:36 Address: ":8080"
2020/10/14 19:44:36 BaseDir: .
2020/10/14 19:44:36 Tokens: [superSecretToken someOtherEquallySuperSecretToken]
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
| `LogFilePath` | string | `""` | The location of the log file where you'd like to log. By default, `gumdrop` will log to `os.Stdout` (STDOUT). | `GUMDROP_LOGFILEPATH` |
| `Tokens` | string array | `[superSecretToken someOtherEquallySuperSecretToken]` | Tokens allowed in `Authorization: bearer Token` header. | _not available_ | 


## Running

`gumdrop` is entirely self-contained, simply run the executable:

```shell script
$ ./gumdrop
2020/10/14 14:56:45 Starting `gumdrop`...
2020/10/14 14:56:45 Address: ":8080"
2020/10/14 14:56:45 BaseDir: .
...
```

## Dropping Files

```shell script
$ echo "some text" > someFile.txt
$ curl -v -X POST \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: bearer superSecretToken" \
  -H "x-directory: someDirectory" \
  -F file=@someFile.txt \
  localhost:8080
$ rm someFile.txt
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