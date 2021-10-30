#!/bin/bash

num_nodes=$1
counter=0
for node in $(cat node_list.txt | head -n $num_nodes); do
    if [[ $counter == 0 ]]
    then
      host=$node
#      echo "$host" > started_nodes.txt
    fi

    echo "$node"
    ssh -f "$node" "cd dht/ && ./dht"
    counter=$((counter+1))
    echo "$node" >> started_nodes.txt
done
