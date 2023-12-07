#!/bin/bash

echo -e "==================="
echo -e "BUILDING NODE IMAGE"
echo -e "===================\n"
# docker build -t node ../.

echo -e "==================="
echo -e "REMOVING OLD NODES"
echo -e "===================\n"

containers=$(docker ps -a -q)
# If the list is not empty, then delete the containers
if [ -n "$containers" ]; then
    echo -e "Deleting existing Docker containers..."
    docker rm -f $containers
else
    echo -e "No Docker containers exist.\n"
fi

# Check if the dynamo network exists and remove it
if [ "$(docker network ls | grep dynamo)" ]; then
  echo -e "\n==================="
  echo -e "REMOVING OLD NETWORK"
  echo -e "==================="
  docker network rm dynamo
fi

echo -e "\n========================"
echo -e "Creating Dynamo Network"
echo -e "========================"

docker network create --driver bridge dynamo

echo -e "\n========================"
echo -e "Running Webclient"
echo -e "========================"

docker run -d -p 8080:8080 --network dynamo --name webclient node ./bin/webclient 

# Define the base port
base_port=50053

# Define the number of nodes
num_nodes=40

echo -e "\n========================"
echo -e "Running nodes 50053, 50054, 50055, 50056, 50051"
echo -e "========================"

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

echo -e "\nWaiting for 10 seconds, letting membership list update"
sleep 10

echo -e "\n========================"
echo -e "PUT foo:bar (Normal use-case)"
echo -e "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar'
echo -e "\n"

echo -e "========================"
echo -e "GET foo expected bar (Normal use-case)"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo -e "\n"

echo -e "========================"
echo -e "Adding new node node-50052"
echo -e "========================"
docker run -d -p 50052:50052 --network dynamo --name "node-50052" node ./bin/server --addr="node-50052:50052" --webclient="http://webclient:8080/addNode?address=node-50052:50052" 

echo -e "Waiting for 10 seconds, letting membership list update\n"
sleep 10

echo -e "========================"
echo -e "PUT foo:bar1 (After adding node-50052)"
echo -e "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar1'
echo -e "\n"

echo -e "========================"
echo -e "PUT foo:bar2 (After adding node-50052)"
echo -e "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar2'
echo -e "\n"

echo -e "========================"
echo -e "PUT foo:bar3 (After adding node-50052)"
echo -e "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar3'
echo -e "\n"

echo -e "========================"
echo -e "GET foo expected bar3 (After adding node-50052)"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo -e "\n"

echo -e "Sleeping for 5 seconds, moving on to next testing phase"
sleep 5

echo -e "========================"
echo -e "Killing node-50052"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/kill' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50052:50052"
}'
echo -e "\n"

echo -e "========================"
echo -e "Killing node-50051"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/kill' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50051:50051"
}'
echo -e "\n"

echo -e "Waiting for 10 seconds, letting membership list update\n"
sleep 10

# echo -e "========================"
# echo -e "GET foo expected bar"
# echo -e "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo -e "\n"

echo -e "========================"
echo -e "PUT foo:hintedhandoff (After killing node-50052 & node-50051)"
echo -e "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=hintedhandoff'
echo -e "\n"

echo -e "========================"
echo -e "GET foo expected hintedhandoff (After killing node-50052 & node-50051)"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo -e "\n"

# echo -e "========================"
# echo -e "GET foo expected hintedhandoff"
# echo -e "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo -e "\n"
#
# echo -e "========================"
# echo -e "GET foo expected hintedhandoff"
# echo -e "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo -e "\n"
#
# echo -e "========================"
# echo -e "GET foo expected hintedhandoff"
# echo -e "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo -e "\n"
#
# echo -e "========================"
# echo -e "GET foo expected hintedhandoff"
# echo -e "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo -e "\n"

echo -e "========================"
echo -e "Reviving node-50052"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/revive' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50052:50052"
}'
echo -e "\n"

echo -e "========================"
echo -e "Reviving node-50051"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/revive' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50051:50051"
}'
echo -e "\n"

echo -e "Waiting for 10 seconds, let hinted handoff do its thing\n"

sleep 10

echo -e "========================"
echo -e "GET foo expected hintedhandoff (Revived node-50052 & node-50051)"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo -e "\n"

echo -e "========================"
echo -e "PUT foo:revivednode1 (Revived node-50052 & node-50051)"
echo -e "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=revivednode1'
echo -e "\n"

echo -e "========================"
echo -e "PUT foo:revivednode2 (Revived node-50052 & node-50051)"
echo -e "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=revivednode3'
echo -e "\n"

echo -e "========================"
echo -e "PUT foo:revivednode3 (Revived node-50052 & node-50051)"
echo -e "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=revivednode3'
echo -e "\n"

echo -e "========================"
echo -e "GET foo expected revivednode3"
echo -e "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo -e "\n"

# echo -e "========================"
# echo -e "GET foo expected hintedhandoff"
# echo -e "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo -e "\n"
#
# echo -e "========================"
# echo -e "GET foo expected hintedhandoff"
# echo -e "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo -e "\n"
#
# echo -e "========================"
# echo -e "GET foo expected hintedhandoff"
# echo -e "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo -e "\n"
