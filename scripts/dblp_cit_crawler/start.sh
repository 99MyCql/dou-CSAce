#!/bin/bash

/bin/rm *.log
nohup go run *.go >log.log &
