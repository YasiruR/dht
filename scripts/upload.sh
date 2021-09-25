#!/bin/bash
echo "building the project"
go build
echo "uploading the binary"
rsync -a dht ywi006@uvcluster.cs.uit.no:/home/ywi006/dht
echo "uploaded successfully"