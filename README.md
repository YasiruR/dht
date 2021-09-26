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

## Tester

A separate testing client is implemented and can be found in tester directory to
test this distributed key-value store. Follow the steps mentioned below to execute 
tester.

1. `cd tester/`
2. `go build`
3. Move this executable file to a node in the cluster
4. `./tester <HTTP_method_type> <host:port> <number_of_requests>`     
  Note: HTTP method type should either be GET or PUT

## Logging

A logger library is integrated for debugging purposes with following 
hierarchical log levels. Required level can be enabled via `log_level` in configs.yaml
by setting the relevant value.

1. ERROR
2. INFO
3. DEBUG
4. TRACE

## Todo

- Implementation of finger tables