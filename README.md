# Distributed Key-Value Store Using Chord

This is the Golang implementation of **chord** with ability to cope with dynamic changes in
the network. 

## Usage

Below shows the steps of using this repository to build and use as a distributed 
key-value store.

1. Build the project using ``go build``
2. Move the executable file (dht) and configs.yaml to relevant nodes in the cluster
3. Update configurations if required
4. Run the executable file

## Tester

### Store Tester

A testing client is implemented and can be found in tester directory to
test distributed key-value store functionality. Follow the steps mentioned below to execute 
tester.

1. `cd tester/store`
2. `go build -o store-tester`
3. Move this executable file to a node in the cluster
4. `./store-tester <HTTP_method_type> <host:port> <number_of_requests>`     
  Note: HTTP method type should either be GET or PUT

### Stability Tester

A separate client is implemented to test network's resilience for dynamic changes in the
structure such as joining, leaving or crashing of nodes. Below steps should be followed for 
the execution of this tester.

1. `cd tester/stability`
2. `go build -o stab-tester`
3. Move this executable file to a node in the cluster
4. Add a text file with corresponding nodes (a sample file is provided) and name it as started_nodes.txt
5. `./stab-tester <join/crash> <number_of_nodes to join or crash>`     
   Note: HTTP method type should either be GET or PUT

Shell scripts can be found in `scripts/cluster` folder for initialization of cluster, joining
and leaving of nodes.

## Logging

A logger library is integrated for debugging purposes with following 
hierarchical log levels. Required level can be enabled via `log_level` in configs.yaml
by setting the relevant value.

1. `ERROR`
2. `INFO`
3. `DEBUG`
4. `TRACE`

## Todo

- Implementation of finger tables