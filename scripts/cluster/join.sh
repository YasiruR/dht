#!/bin/bash

counter=0
for node in $(cat started_nodes.txt); do
    # shellcheck disable=SC1073
    if [[ $counter == 0 ]]
    then
      host=$node
      counter=$((counter+1))
      continue
    fi
    curl -X POST --retry 3 --retry-delay 0 "${node}:52520/join?nprime=${host}:52520"
    counter=$((counter+1))
done
