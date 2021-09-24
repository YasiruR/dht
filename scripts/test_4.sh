#!/bin/bash

# this test is carried out for node-2 (compute-3-28), node-5 (compute-8-7), node-8 (compute-6-17) and node-11 (compute-6-22)
curl -X POST -H "Content-Type: text/plain" --data "value-999" compute-3-28:52520/storage/999
if curl -X GET -H "Content-Type: text/plain" compute-3-28:52520/storage/999 | grep -q 'value-999'; then
  echo "key=2 ok"
else
  echo "key=2 not ok"
fi

curl -X POST -H "Content-Type: text/plain" --data "value-12345" compute-3-28:52520/storage/12345
if curl -X GET -H "Content-Type: text/plain" compute-3-28:52520/storage/12345 | grep -q 'value-12345'; then
  echo "key=5 ok"
else
  echo "key=5 not ok"
fi

if curl -X GET -H "Content-Type: text/plain" compute-8-7:52520/storage/12345 | grep -q 'value-12345'; then
  echo "key=5 ok 2"
else
  echo "key=5 not ok 2"
fi

curl -X POST -H "Content-Type: text/plain" --data "value-abcdef" compute-8-7:52520/storage/abcdef
if curl -X GET -H "Content-Type: text/plain" compute-6-17:52520/storage/abcdef | grep -q 'value-abcdef'; then
  echo "key=1 ok"
else
  echo "key=1 not ok"
fi

curl -X POST -H "Content-Type: text/plain" --data "value-abc123" compute-6-17:52520/storage/abc123
if curl -X GET -H "Content-Type: text/plain" compute-6-22:52520/storage/abc123 | grep -q 'value-abc123'; then
  echo "key=0 ok"
else
  echo "key=0 not ok"
fi

curl -X POST -H "Content-Type: text/plain" --data "value-assignment" compute-6-22:52520/storage/assignment
if curl -X GET -H "Content-Type: text/plain" compute-3-28:52520/storage/assignment | grep -q 'value-assignment'; then
  echo "key=1 ok"
else
  echo "key=1 not ok"
fi

curl -X POST -H "Content-Type: text/plain" --data "value-dht" compute-8-7:52520/storage/dht
if curl -X GET -H "Content-Type: text/plain" compute-6-17:52520/storage/dht | grep -q 'value-dht'; then
  echo "key=0 ok"
else
  echo "key=0 not ok"
fi

