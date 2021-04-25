#!/bin/bash

ps -aux | grep main.go | grep -v grep | awk '{print $2}' | xargs kill -9
ps -aux | grep auto_exec.sh | grep -v grep | awk '{print $2}' | xargs kill -9
ps -aux | grep exe/main | grep -v grep | awk '{print $2}' | xargs kill -9
