# Distributed Key-Value Store Using Chord

This is the Golang implementation of **chord**.

## Usage

Below shows the steps of using this repository to build and use as a distributed 
key-value store.

1. Build the project using ``go build``
2. Move the executable file (dht) and configs.yaml to relevant nodes in the cluster
3. Update predecessor and successor hostnames in configs.yaml as follows:
```yaml
# neighbour configs
finger_table_enabled: false
neighbour_check: false
predecessor: "compute-1-1"
predecessor_port: ""
successor: "compute-2-1"
successor_port: ""
```
NOTE: Ports can be given null if they are same as the corresponding node
4. Run the executable file

## Todo

- Implementation of finger tables