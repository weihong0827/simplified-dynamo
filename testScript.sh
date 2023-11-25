#!/bin/bash

echo "==================="
echo "BUILDING NODE IMAGE"
echo "==================="
docker build -t node .

echo ""

echo "==================="
echo "REMOVING OLD NODES"
echo "==================="
containers = $(docker ps -a -q)
if [ -n "$containers" ]; then
    echo "Deleting existing Docker containers..."
    docker rm -f $containers
else
    echo "No Docker containers exist."
fi
echo ""

# Check if the dynamo network exists and remove it
if [ "$(docker network ls | grep dynamo)" ]; then
  echo "==================="
  echo "REMOVING OLD NETWORK"
  echo "==================="
  docker network rm dynamo
  echo ""
fi

echo "========================"
echo "Creating Dynamo Network"
echo "========================"

docker network create --driver bridge dynamo

echo ""

echo "========================"
echo "Running Webclient"
echo "========================"

docker run -d -p 8080:8080 --network dynamo --name webclient node ./bin/webclient 

# Define the base port
base_port=50053

# Define the number of nodes
num_nodes=4

echo ""

echo "========================"
echo "Running nodes 50053, 50054, 50055, 50056, 50051"
echo "========================"

# Loop to create nodes
for (( i=0; i<num_nodes; i++ )); do
  # Calculate the current port
  port=$((base_port+i))

  # Generate the node name
  node_name="node-$port"

  # Run the Docker command
  docker run -d -p $port:$port --network dynamo --name "$node_name" node ./bin/server --addr="$node_name:$port" --webclient="http://webclient:8080/addNode?address=$node_name:$port"
done

docker run -d -p 50051:50051 --network dynamo --name "node-50051" node ./bin/server --addr="node-50051:50051" --webclient="http://webclient:8080/addNode?address=node-50051:50051" 

echo ""

sleep 12

echo "========================"
echo "PUT foo:bar"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar'
echo ""
echo ""

echo "========================"
echo "GET foo expected bar"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""
echo ""

sleep 3

echo "========================"
echo "Adding new node node-50052"
echo "========================"
docker run -d -p 50052:50052 --network dynamo --name "node-50052" node ./bin/server --addr="node-50052:50052" --webclient="http://webclient:8080/addNode?address=node-50052:50052" 
echo ""

sleep 3

echo "========================"
echo "GET foo expected bar"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""
echo ""

## Killing
echo "========================"
echo 'Killing node-50052'
echo "========================"
curl --location 'http://127.0.0.1:8080/kill' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50052:50052"
}'
echo ""
echo ""

# echo "========================"
# echo 'Killing node-50051'
# echo "========================"
# curl --location 'http://127.0.0.1:8080/kill' \
# --header 'Content-Type: application/json' \
# --data '{
#     "Address": "node-50051:50051"
# }'
# echo ""
# echo ""

sleep 5

echo "========================"
echo 'GET foo expected bar'
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""
echo ""

echo "========================"
echo "PUT foo:kill50052"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=kill50052'
echo ""
echo ""

echo "========================"
echo "GET foo expected kill50052"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""
echo ""

## Reviving
echo "========================"
echo "Reviving node-50052"
echo "========================"
curl --location 'http://127.0.0.1:8080/revive' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50052:50052"
}'
echo ""
echo ""

echo "========================"
echo "GET foo expected kill2nodes"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""
echo ""

echo "========================"
echo "PUT foo:notbar"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=notbar'
echo ""
echo ""

echo "========================"
echo "GET foo expected notbar"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""
echo ""

echo "========================"
echo "PUT foo:notbar2"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=notbar2'
echo ""
echo ""

echo "========================"
echo "PUT foo:notbar3"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=notbar3'
echo ""
echo ""

# echo "========================"
# echo "PUT foo:notbar4"
# echo "========================"
# curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=notbar'
# echo ""
# echo ""

# echo "========================"
# echo "PUT foo:notbar5"
# echo "========================"
# curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=notbar'
# echo ""
# echo ""
#
# echo "========================"
# echo "PUT foo:hello"
# echo "========================"
# curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=hello'
# echo ""
# echo ""

echo "========================"
echo "GET foo expected notbar3"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""
echo ""
