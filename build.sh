#/bin/bash

go build -trimpath -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -o bin/zrouter
