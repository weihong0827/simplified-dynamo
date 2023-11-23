#!/bin/bash

docker build -t node .

docker rm -f $(docker ps -a -q)

# docker network create --driver bridge dynamo

docker run -d -p 8080:8080 --network dynamo --name webclient node ./bin/webclient 

# Define the base port
base_port=50051

# Define the number of nodes
num_nodes=5

# Loop to create nodes
for (( i=0; i<num_nodes; i++ )); do
  # Calculate the current port
  port=$((base_port+i))

  # Generate the node name
  node_name="node-$port"

  # Run the Docker command
  docker run -d -p $port:$port --network dynamo --name "$node_name" node ./bin/server --addr="$node_name:$port" --webclient="http://webclient:8080/addNode?address=$node_name:$port"
done

go run test/test.go # GET and PUT

additionalNode=$((base_port + num_nodes))

docker run -d -p $additionalNode:$additionalNode --network dynamo --name "node-$additionalNode" node ./bin/server --addr="node-$additionalNode:$additionalNode" --webclient="http://webclient:8080/addNode?address=node-$additionalNode:$additionalNode"

# go run test2.go # GET and PUT with failure
