#!/bin/sh

echo "BUILDING NODE IMAGE"
echo "==================="
docker build -t node .

docker rm -f $(docker ps -a -q)

# Check if the dynamo network exists and remove it
if [ "$(docker network ls | grep dynamo)" ]; then
    docker network rm dynamo
fi

docker network create --driver bridge dynamo

echo "Running Webclient"

docker run -d -p 8080:8080 --network dynamo --name webclient node ./bin/webclient 

# Define the base port
base_port=50051

# Define the number of nodes
num_nodes=$1

echo "Running $num_nodes nodes"

# Loop to create nodes
for (( i=0; i<num_nodes; i++ )); do
  # Calculate the current port
  port=$((base_port+i))

  # Generate the node name
  node_name="node-$port"

  # Run the Docker command
  docker run -d -p $port:$port --network dynamo --name "$node_name" node ./bin/server --addr="$node_name:$port" --webclient="http://webclient:8080/addNode?address=$node_name:$port"
done