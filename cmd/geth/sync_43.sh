#!/bin/bash

sshpass -p "123456" scp geth eth@192.168.136.43:~/gm/
sshpass -p "123456" scp genesis.json eth@192.168.136.43:~/gm/