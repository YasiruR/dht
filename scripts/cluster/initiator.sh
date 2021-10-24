#!/bin/bash

num_nodes=$1
counter=0
for node in $(cat node_list.txt | head -n $num_nodes); do
    if [[ $counter == 0 ]]
    then
      host=$node
    fi

    echo "$node"
    ssh -f "$node" "cd dht/ && ./dht"
    counter=$((counter+1))
    echo "$node" >> started_nodes.txt
done

## shellcheck disable=SC2028
#echo "$host" > started_nodes.txt
#for node in $(cat node_list.txt | head -n $num_nodes); do
#    # shellcheck disable=SC1073
#    if [ "$node" == "$host" ]
#    then
#      continue
#    fi
#    curl -X POST --retry 3 --retry-delay 0 "${node}:52520/join?nprime=${host}:52520"
#    # shellcheck disable=SC2028
#    echo "$node" >> started_nodes.txt
#done