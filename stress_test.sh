#!/bin/bash
for i in {1..10}
do
    go run cmd/client/main.go testdata/sample.txt &
done
wait 
echo "All uploads sent!"