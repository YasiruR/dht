#!/bin/bash

num_leaves=$1
counter=0
for node in $(cat started_nodes.txt); do
    # shellcheck disable=SC2053
    if [[ $counter == $num_leaves ]]
    then
        break
    fi
    curl -X POST --retry 3 --retry-delay 0 "${node}:52520/leave"
    counter=$((counter+1))
done
