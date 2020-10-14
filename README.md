# gumdrop

## About

`gumdrop` is meant to provide a simple HTTP interface for dropping files onto a machine.  It's basically an FTP `PUT` operation when you don't have or want FTP.

## Building

NOTE: This is a `go` application.  Please have `go` >= 1.15 installed before building.

```shell script
$ git clone git@github.com:charlesread/gumdrop.git
$ cd gumdrop
$ make build
```

This will build `bin/gumdrop`.

## Installing

To build and install `gumdrop` at `/usr/local/bin/gumdrop` simply:

```shell script
$ make install
``` 

### Install Service

`gumdrop` runs just fine from the command line, a `systemd` unit file is provided at `gumdrop.service`.

The default service runs as the user `gumdrop`, be sure that user exists:

```shell script
$ useradd gumdrop -s /sbin/nologin -m
$ cp config.yaml /home/gumdrop # copy over default config for convenience (optional)
$ chown gumdrop:gumdrop /home/gumdrop/config.yaml
```

You can install the service with the `service` target:

```shell script
$ make install
$ make service
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

| Name | Default Value | Function |
| ---- | ------------- | -------- |
| `Address` | `:8080` | Sets the address where `gumdrop` will serve. |
| `BaseDir` | `/opt/misc/drops` | The base directory where files will be dropped. |


## Running

`gumdrop` is entirely self-contained, simply run the executable:

```shell script
$ ./gumdrop
2020/10/14 14:56:45 Starting `gumdrop`...
2020/10/14 14:56:45 Address: ":8090"
2020/10/14 14:56:45 BaseDir: /opt/misc/drops
...
```

## Dropping Files

```shell script
$ curl -v -X POST \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: bearer abc123" \
  -H "x-directory: someDirectory" \
  -F file=@beep.txt \
  localhost:8080/drop
```