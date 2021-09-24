#!/bin/bash
echo "building the project"
go build
echo "uploading the binary"
rsync -a dht ywi006@uvcluster.cs.uit.no:/home/ywi006/dht
#echo "logging into uv cluster"
#ssh ywi006@uvcluster.cs.uit.no
#echo "copying binary to respective node directories"
#cp /home/ywi006/dht/dht /home/ywi006/dht/node-2/
#cp /home/ywi006/dht/dht /home/ywi006/dht/node-5/
#cp /home/ywi006/dht/dht /home/ywi006/dht/node-8/
#cp /home/ywi006/dht/dht /home/ywi006/dht/node-11/
#echo "starting nodes"
#ssh compute-3-28
#/home/ywi006/dht/node-2/dht &
#ssh compute-8-7
#/home/ywi006/dht/node-5/dht &
#ssh compute-6-17
#/home/ywi006/dht/node-8/dht &
#ssh compute-6-22
#/home/ywi006/dht/node-11/dht &
echo "uploaded successfully"
