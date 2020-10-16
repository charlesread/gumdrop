#!/usr/bin/env bash

cd /tmp
git clone https://github.com/charlesread/gumdrop.git
cd gumdrop
useradd gumdrop -s /sbin/nologin -m
cp config.yaml /home/gumdrop
chown gumdrop:gumdrop /home/gumdrop/config.yaml # edit appropriately
make install
make service