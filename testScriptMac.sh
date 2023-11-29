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
base_port=50053

# Define the number of nodes
num_nodes=4

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

docker run -d -p 50051:50051 --network dynamo --name "node-50051" node ./bin/server --addr="node-50051:50051" --webclient="http://webclient:8080/addNode?address=node-50051:50051" 

echo "========================"

sleep 5

echo "PUT foo:bar"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar'
echo ""

echo "GET foo"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""

echo "PUT foo:bar"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar'
echo ""

echo "GET foo"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""

sleep 3

echo "Adding new node node-50052"
echo "========================"
docker run -d -p 50052:50052 --network dynamo --name "node-50052" node ./bin/server --addr="node-50052:50052" --webclient="http://webclient:8080/addNode?address=node-50052:50052" 

sleep 5

echo "Get foo"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""

echo "PUT foo:bar"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar'
echo ""

echo "GET foo"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""

echo "PUT foo:bar"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=bar'
echo ""

echo "GET foo"
echo "========================"
curl --location 'http://127.0.0.1:8080/get?key=foo'
echo ""

echo 'Killing node-50052'
echo "========================"
curl --location 'http://127.0.0.1:8080/kill' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50052:50052"
}'
echo ""

echo 'Killing node-50051'
echo "========================"
curl --location 'http://127.0.0.1:8080/kill' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50051:50051"
}'
echo ""

sleep 5


echo "Put foo:notbar"
echo "========================"
curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=notbar'
echo ""

# echo 'Get "foo"'
# echo "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo ""

# echo "Reviving node-50052"
# echo "========================"
# curl --location 'http://127.0.0.1:8080/revive' \
# --header 'Content-Type: application/json' \
# --data '{
#     "Address": "node-50052:50052"
# }'
# echo ""

# echo "Get foo"
# echo "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
# echo ""

# echo "Put foo:notbar"
# echo "========================"
# curl --location --request PUT 'http://127.0.0.1:8080/put?key=foo&value=notbar'
# echo ""

# echo "Get foo"
# echo "========================"
# curl --location 'http://127.0.0.1:8080/get?key=foo'
