#!/bin/bash

set -ex

go build --tags sm2 -a -ldflags '-extldflags "-static"' .

sshpass -p "123456" scp geth eth@192.168.136.43:~/gm/
sshpass -p "123456" scp genesis.json eth@192.168.136.43:~/gm/