#!/bin/bash
echo "building the project"
cd ..
go build
echo "uploading the binary"
cd scripts/
rsync -a ../dht ywi006@uvcluster.cs.uit.no:/home/ywi006/dht
echo "uploaded successfully"