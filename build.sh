#/bi/bash

go build -trimpath -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH
