#!/bin/bash

echo "export PATH=$PATH:/usr/local/go/bin" >>~/.profile
source ~/.profile

tmux new-session -d -s api './start-peer2.sh'
sleep 60
ls
tmux detach -s api

echo "API started in detached tmux session."