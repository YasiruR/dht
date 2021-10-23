#!/bin/bash

num_nodes=$1
for node in $(cat node_list.txt | head -n $num_nodes); do
    echo "$node"
    ssh -f "$node" "cd dht/ && ./dht"
done