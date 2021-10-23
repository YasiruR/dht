#!/bin/bash

num_nodes=$1
host=$2
for node in $(cat node_list.txt | head -n $num_nodes); do
    echo "$node"
    ssh -f "$node" "cd dht/ && ./dht"
done

for node in $(cat node_list.txt | head -n $num_nodes); do
    # shellcheck disable=SC1073
    if [ "$node" == "$host" ]
    then
      continue
    fi
    curl -X POST "${node}:52520/join?nprime=${host}:52520"
done