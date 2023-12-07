<div align="center">
<h1 align="center">
<img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-markdown-open.svg" width="100" />
<br></h1>
<h3>simplified-dynamo</h3>

<p align="center">
<img src="https://img.shields.io/badge/GNU%20Bash-4EAA25.svg?style=flat-square&logo=GNU-Bash&logoColor=white" alt="GNU%20Bash" />
<img src="https://img.shields.io/badge/YAML-CB171E.svg?style=flat-square&logo=YAML&logoColor=white" alt="YAML" />
<img src="https://img.shields.io/badge/Docker-2496ED.svg?style=flat-square&logo=Docker&logoColor=white" alt="Docker" />
<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat-square&logo=Go&logoColor=white" alt="Go" />
</p>
</div>

---

## 📖 Table of Contents
- [📖 Table of Contents](#-table-of-contents)
- [📍 Overview](#-overview)
- [📦 Features](#-features)
- [📂 Repository Structure](#-repository-structure)
- [⚙️ Modules](#️-modules)
- [🚀 Getting Started](#-getting-started)
  - [🔧 Installation](#-installation)
  - [🤖 Running](#-running)
  - [Running the script](#running-the-script)
  - [🧪 Tests](#-tests)
- [Contributors ✨](#contributors-)

## 📍 Overview

The repository houses a Go-based distributed system with a key-value store that uses Docker for containerization and gRPC for remote procedure calls. It utilizes vector clocks for concurrent control, hash functions for data distribution, and manages servers with custom scripts. The system demonstrates fault tolerance through functionalities like hinted handoff, node revival, and operation success upon pre-defined responses. Users can interact via a web client while automatically building and testing services aided by Docker and Makefile. This facilitates robust, highly available, scalable, and distributed data management.

[See our full report here.](https://docs.google.com/document/d/1nZMDEix41mRh7ARzVUOCA1t1cr6zjSowzWhDP2hK2ns/edit?usp=sharing)

## 📦 Features

Just simply run `docker-compose up --build` in the terminal and it should build the network for you

**GET**: `http://127.0.0.1:8080/get?key=foo`

**PUT**: `http://127.0.0.1:8080/put?key=foo&value=bar`

The desired key and value are input as queries into the URL as shown above. For easy testing, you may use a system such as Postman to assist in making the API calls.



## 📂 Repository Structure

```sh
└── /
    ├── Dockerfile
    ├── Key-range.txt
    ├── Makefile
    ├── buildNodesMac.sh
    ├── config/
    │   ├── config.go
    │   └── enums.go
    ├── docker-compose.yml
    ├── go.mod
    ├── go.sum
    ├── hash/
    │   ├── error.go
    │   ├── hash.go
    │   └── hash_test.go
    ├── main.go
    ├── pb/
    │   ├── dynamo.pb.go
    │   └── dynamo_grpc.pb.go
    ├── proto/
    │   └── dynamo.proto
    ├── scripts/
    │   ├── buildScript.sh
    │   ├── hintedHandoff.sh
    │   ├── testScript.sh
    │   └── testScriptMac.sh
    ├── server/
    │   ├── error.go
    │   ├── grpc.go
    │   ├── keyrange.go
    │   ├── main.go
    │   ├── ops.go
    │   ├── utils.go
    │   ├── utils_test.go
    │   ├── vectorClock.go
    │   └── vectorClock_test.go
    ├── test/
    │   └── load_test.go
    └── webclient/
        └── main.go

```

## ⚙️ Modules

<details closed><summary>Root</summary>

| File                              | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| ---                               | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| [go.mod]({file_path})             | The project uses Go for building a distributed system with gRPC for communication. It comprises server operations, hashing features, configuration, protocol buffers definitions and corresponding gRPC stubs, scripts for build/tests, Docker configurations for containerization, and a web client. It leverages several dependencies like gin for web functionality, grpc for RPC support, and protobuf for data serialization. Unit tests and load tests are present for quality assurance.            |
| [Dockerfile]({file_path})         | The code is for a Dockerized Go application with gRPC functionalities. It includes modules for configurations, hashing, protobuf, server operations, and web client. Various scripts for building and testing are present. The Dockerfile builds the server and web client from the latest Go image, downloads dependencies with Go Mod, and exposes port 8080.                                                                                                                                            |
| [Makefile]({file_path})           | The code includes scripts to compile protocol buffers specification files, run a web client, initiate and run five server nodes (each on a different port), all locally. It also contains instructions to build a Docker image tagged as node, remove and rerun two specific server nodes (second and sixth) within Docker-one operating on the host network and the other on a network named dynamo_default.                                                                                              |
| [Key-range.txt]({file_path})      | The code outlined above is for a distributed hash table application built with Go, gRPC, and Docker. It consists of a Dockerfile and scripts for building and testing, configuration files, protobuffers files for implementing gRPC, and specific Go files for managing hashing, server operations, and concurrency control with vector clocks. The Key-range.txt file likely drives key distribution across the network nodes. The Docker and Makefile automate the environment setup and build process. |
| [go.sum]({file_path})             | This code includes a directory tree that is structured for a Go-based project with Docker. It provides functionalities for hash handling, config generation, protocol buffer interactions, server and web-client operations. Scripts for building and testing as well as Docker and Makefile for containerization and automation are present. go.sum includes dependencies for the project like gin-gonic, go-playground, etc.                                                                             |
| [buildNodesMac.sh]({file_path})   | The code is intended to automate the process of setting up a Docker-based distributed system. It first builds a Docker image named node and deletes any existing containers and networks. Following this, it creates a new network named dynamo and starts running a web client container in this network. The script then brings up a user-specified number of node containers at sequential ports in the dynamo network. Each node also registers its address with the web client.                       |
| [docker-compose.yml]({file_path}) | The provided Docker compose YAML file facilitates managing a set of microservices in their separate containers. It specifies six services including a web client and five servers (also known as nodes). The web client and nodes communicate via specified ports. Each server depends on the web client being healthy before starting. The servers also register themselves with the web client upon startup. The configuration allows them to operate in a shared dynamo network.                        |
| [main.go]({file_path})            | This codebase facilitates the creation of a distributed system using Go, gRPC protocols, and Docker. The architecture includes configuration entities, hashing utilities, server-side operations, protofiles for gRPC service definition, and test scripts. It also involves key-range handling, vector clock for maintaining synchronization among nodes, and error handling. main.go functions are entry points across modules while Dockerfile and docker-compose.yml handle containerization.          |

</details>

<details closed><summary>Webclient</summary>

| File                   | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| ---                    | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| [main.go]({file_path}) | The provided code is for a distributed systems client application built in Go that uses gRPC for communication. Primary functions include creating connections to multiple servers, adding nodes, killing and reviving nodes, and reading and writing operations. It finds the fastest responding server to handle GET and PUT requests, and supports operations like reading and writing to a key-value store, with responses converted to a readable format. It also provides health checks and facilities for adding, killing and reviving server nodes. All these operations are exposed as HTTP API endpoints. |

</details>

<details closed><summary>Pb</summary>

| File                             | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| ---                              | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| [dynamo_grpc.pb.go]({file_path}) | This code defines an API interface for a GRPC-based key-value store service named dynamo. The API has functions for different operations such as Write, Read, Join, Gossip, Delete, and KillNode. It also includes methods for handling hinted hand-off reads/writes and bulk writes. The code embeds UnimplementedKeyValueStoreServer for forward compatibility and defines the GRPC service description for the KeyValueStore service.                                                                               |
| [dynamo.pb.go]({file_path})      | The code is a structured Go project, equipped with Docker configurations for containerization and gRPC protocol for client-server communication. It utilises Hashing and Vector Clocks for distributed systems operations. Test files and scripts are available to validate functionality. The displayed file, dynamo.pb.go, is auto-generated using Protocol Buffers for serialization of structured data in the dynamo application. It includes packages related to Protobuf reflection and timestamp specification. |

</details>

<details closed><summary>Proto</summary>

| File                        | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| ---                         | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| [dynamo.proto]({file_path}) | The code outlines a distributed key-value store built using Google's Protocol Buffers (protobuf). It supports operations like Write, Read, Forward, Join, Gossip, and various Replication-handling like HintedHandoff and Bulk operations. Node management operations include KillNode and ReviveNode. It uses vector clocks for concurrency control and membership lists for managing nodes in the system. The proto file illustrates structure for message interchange, with requests and responses for various operations. A Docker setup is indicated with Dockerfile and docker-compose.yml for containerized deployment. |

</details>

<details closed><summary>Test</summary>

| File                        | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| ---                         | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| [load_test.go]({file_path}) | The code performs load tests on PUT and GET HTTP requests. It uses Go's testing package to concurrently send 10 put and get requests to a local server using Goroutines and waits for all requests to finish using a WaitGroup. It logs the time taken for all requests to complete and checks whether response status is OK or not. The tests simulate loading conditions on the server and measure the performance of the system under such a load. |

</details>

<details closed><summary>Hash</summary>

| File                        | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| ---                         | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| [error.go]({file_path})     | The displayed code structure represents a distributed system based on gRPC with Docker containerization. Key components include configuration files, a hashing function, protocol buffers for gRPC, server-side logic including error handling and vector clock management, automated test scripts, and a dedicated web client. A specific error handling mechanism is highlighted in the hash/error.go code snippet, defining standard error messages.                                          |
| [hash_test.go]({file_path}) | The code contains several test cases that validate the functionality of GetAddressFromNode and GetNodesFromKey functions in different scenarios. These include checks for an empty node slice, when a single node or multiple nodes are present, if a key results in a hash out of the nodes' range, and when the nodes are unsorted. The tests assert that the correct node is retrieved or error messages are appropriately returned.                                                          |
| [hash.go]({file_path})      | The code is part of a Go project which implements a distributed hash table (DHT) using the Dynamo-style data partitioning and replication. It defines helper functions to find nodes, based on a given key, in a sorted nodes list. It also has a function to generate a hash from a given key. The identified nodes are then used for data storage or retrieval operations, facilitating distributed data storage with fault tolerance. The binary search operation ensures quick node lookups. |

</details>

<details closed><summary>Config</summary>

| File                     | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| ---                      | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| [config.go]({file_path}) | The directory tree visualizes the structure of a Go-based distributed system. It uses Docker for containerization, allows code testing, employs Protocol Buffers for serializing structured data, and uses gRPC for remote procedure calls. The file `config.go` under the `config` directory sets the values of read (R), write (W) and node (N) operations for the system. The system is partitioned into subdirectories each hosting different functionalities, such as hashing, server operations, and tests. |
| [enums.go]({file_path})  | The code defines configuration constants for a distributed system. It includes different types of `Operation` (READ, WRITE) and `KeyRangeOp` (TRANSFER, DELETE). There is a function to convert these operations to their string equivalents. It's part of a larger system, containing server logic, hashing methods, protobuf files, testing scripts, Docker setup, and a web client.                                                                                                                            |

</details>

<details closed><summary>Server</summary>

| File                               | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| ---                                | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| [vectorClock.go]({file_path})      | The provided Go code in the `vectorClock.go` file is part of a larger project, and it mainly deals with comparing and managing vector clocks in a distributed system. It comprises a function to compare vector clocks of key-value pairs, identify concurrent pairs, update or delete data based on the comparison results, and a utility function to delete a specific element from the data array. The code aims to ensure synchronization and consistency in key-value stores across the system.                                                                                                                            |
| [error.go]({file_path})            | The code is part of a larger Go-based project structure. It defines two global error variables in the server/error.go file, namely NodeDead and NodeALive, to standardize error messages related to the status of a node in a network. These errors can be leveraged throughout the codebase for consistent error handling and reporting.                                                                                                                                                                                                                                                                                       |
| [utils_test.go]({file_path})       | This codebase represents a containerized Go application setup. The core functionalities include configuration handling, hash management, server operations and GRPC services. A series of shell scripts and makefile are provided for building the app. In the given `utils_test.go` file, two test cases are defined to validate the `IsKeyInRange` function: one checks if the function correctly identifies a key in a range and the other confirms that it accurately recognizes when a key falls outside the range.                                                                                                        |
| [grpc.go]({file_path})             | This project revolves around a GRPC-based system built in Go, utilizing Docker for containerization. The code orchestrates a variety of functionalities such as key-range handling, vector clocks, error handling, and hashing. It also includes testing and script files for system building and system testing. Specifically, the grpc.go file contains a function to establish an insecure GRPC connection to a given address, returning any connection errors if they occur.                                                                                                                                                |
| [keyrange.go]({file_path})         | The code involves key-value pair management in a distributed storage system via gRPC. Functions Transfer and BulkWriteToTarget move key-value pairs within specific ranges to a target node while DeleteReplicaFromTarget deletes replica from the target. This data movement is essential during node addition and deletion, load balancing, and failure recovery in the storage system. Transfer logs errors upon occurrence and retries data transfer if necessary.                                                                                                                                                          |
| [utils.go]({file_path})            | The provided codebase is for a GO-based distributed system run in a Docker container. It includes a GRPC server, handling operations on key ranges, vector clocks, and errors. Hash functions, tests, and protocols are also outlined. Separate scripts are constructed to build and test the system. The `IsKeyInRange` function, part of the server utilities, checks if a given key lies within a specified range.                                                                                                                                                                                                           |
| [vectorClock_test.go]({file_path}) | The provided code is primarily a collection of Go unit tests for a vector clock, utilized to manage data consistency in distributed systems. It defines various scenarios checking if the clock is concurrent, ahead, behind, or equal to another. Additionally, it tests vector clock comparisons for handling different data discrepancies during read operations.                                                                                                                                                                                                                                                            |
| [ops.go]({file_path})              | The code comprises multiple gRPC operations for a distributed key-value store. It contains functions for creating gRPC calls, forming connections to nodes, triggering operations, and handling errors. Particular operations include reading and writing key-value pairs, along with a hinted handoff mechanism for error tolerance. Responses are gathered using goroutines and channels, employing context handling and cancellation. Operations are encapsulated within a context-based time-out setting to manage failures. Callers wait for a predefined number of responses before considering the operation successful. |
| [main.go]({file_path})             | The code implements a distributed key-value store, following the DynamoDB system's principles. This hash-based system leverages gRPC for intra-cluster communication and provides functionalities like joins, hash-based key-value storage, membership list reconciliation via gossip protocol, hinted handoff mechanism for temporary failure handling that allows writes to return successfully, and vector clocks for resolving read conflicts among different replicas. All operations are thread-safe with lock protection. It also periodically checks for dead nodes and revives them when applicable.                   |

</details>

<details closed><summary>Scripts</summary>

| File                            | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| ---                             | ---                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| [buildScript.sh]({file_path})   | The bash script manages a Docker-based Dynamo-style distributed database: building a Docker image named node, deleting existing Docker containers, and creating a Docker network named dynamo. It then launches multiple nodes and a web client inside the network. Requests to store a key-value pair foo:bar are made, a new node is added, and the key is verified. Node-50052 is then killed and later revived, testing for fault tolerance. All actions are performed with the aid of curl commands.                                                                   |
| [testScript.sh]({file_path})    | The script builds a Docker image, removes existing nodes and network, then creates a dynamo network and webclient. It runs multiple nodes, does PUT and GET operations for the key foo, adds a new node, kills and revives it, and records the GET results. The nodes communicate with the webclient via HTTP commands, enabling data manipulation and maintaining system operations even when a node is killed.                                                                                                                                                            |
| [testScriptMac.sh]({file_path}) | The script initializes a Docker network, runs a web client and multiple server nodes on it. Using HTTP requests, the script puts pairs of keys and values into the system, retrieves them, and prints the result. It also interacts with the network by adding new nodes, killing existing nodes, and reviving them while observing how these operations impact data retrieval. Additionally, it tidies up any previously existing Docker images and networks at the start.                                                                                                 |
| [hintedHandoff.sh]({file_path}) | The Bash script `hintedHandoff.sh` sets up a networked environment using Docker, deploys a web client and multiple instances of a server, and performs a series of operations that mimic a distributed system undergoing changes. This includes creating nodes, sending PUT and GET requests to store and retrieve data respectively, adding new nodes, intentionally causing node failures, and reviving failed nodes. It showcases the concept of hintedhandoff-a data replication technique that ensures data availability in a distributed system during node failures. |

</details>

---

## 🚀 Getting Started

***Dependencies***

Please ensure you have the following dependencies installed on your system:

* Docker

We require docker as we have containerised or nodes and applications to allow for easier deployment across all machines.

You can also play with `docker compose` file to spin up the network and the `Makefile`'s `run-<num>-node` command to add new nodes to the network

### 🔧 Installation
[Docker Installation Guide](https://docs.docker.com/engine/install/)

### 🤖 Running 

Now you can access the webclient on [localhost:8080](http://localhost:8080) You can interact with the database and the node clusters through the following endpoints:
1. **/get?key=[keyToGet]**: This endpoint retrieve the kv pair from the database and return it as the following format:
```
{
    "message": {
        "Key": "foo",
        "ReplicaResponse": [
            {
                "Value": "bar3",
                "VectorClock": {
                    "429502189": 1
                }
            }
        ],
        "Hashvalue": 2898073819
    }
}
```
Where the `ReplicaResponse` is the result from each of the responsible nodes, if there isnt a consistency problem, only one response will be returned, else more than 1 response will return
2. **/put?key=[keyToPut]&value=[valueToPut]**
You will get a response of `Write Successfully` if there is at least `W` nodes perform a successful write
3. **/kill**
send a `POST` request to this endpoint with the body
```
{
   Address:<node name>:<node-port>
}
```
```bash
## sample request
curl --location 'http://127.0.0.1:8080/kill' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50052:50052"
}'
```
4. **/revive**
send a `POST` request to this endpoint with the body
```
{
   Address:<node name>:<node-port>
}
```
```bash
## sample request
curl --location 'http://127.0.0.1:8080/revive' \
--header 'Content-Type: application/json' \
--data '{
    "Address": "node-50052:50052"
}'
```
### Running the script
Make sure you have docker installed and you are on a UNIX base system

go to the `script` directory, and use the different scripts to test out the different functionalities 

### 🧪 Tests
Ensure that your **docker daemon** is running before running the tests.

To perform **correctness testing** as shown in the report, go into the `scripts` folder and run the `./presentation.sh` file. 

```sh
cd ./scripts
./presentation.sh
```

To perform **load testing** as shown in the report, ensure you are in the `scripts` folder and run `./loadTest.sh`

```sh
cd ./scripts
./loadTest.sh
```
---

## Contributors ✨

Thanks go to these wonderful people:

<a href="https://github.com/weihong0827/simplified-dynamo/graphs/contributors">
<img src="https://contrib.rocks/image?repo=weihong0827/simplified-dynamo" />
</a>
