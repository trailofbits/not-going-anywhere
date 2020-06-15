#!/bin/sh
#@(#) just a simple example of how to build to ProtoBuf definition 
#@(#) We used for Not Going Anywhere

PPATH=$GOPATH/src/github.com/trailofbits/not-going-anywhere/internal/friends
echo $PPATH

if [ -d "$PPATH" ]
then
    protoc --proto_path=$PPATH --go-grpc_out=$PPATH --go_out=$PPATH --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative $PPATH/friends.proto
else 
    echo "could not find directory"
fi
