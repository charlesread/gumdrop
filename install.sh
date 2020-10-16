#!/usr/bin/env bash

useradd gumdrop -s /sbin/nologin -m
cp config.yaml /home/gumdrop
chown gumdrop:gumdrop /home/gumdrop/config.yaml # edit appropriately
make install
make service