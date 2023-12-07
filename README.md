# Overview

We will be implementing a simple clone of `AWS dynamo`

This project implements a distributed key-value store with the data arranged in a ring structure for load balancing and fault tolerance. The system ensures data consistency and availability through replication, consistent hashing, and a gossip-based protocol for service discovery.

The specific breakdown and understanding of dynamo can be found in the blog [here](https://blog.weihong.tech/Projects/Distributed-Systems/Introduction)

# Features

- **Ring Structure Arrangement:** The nodes in the system are arranged in a ring structure using consistent hashing. This setup aids in load distribution and locating key-value pairs efficiently.
- **Replication:** Each data item is replicated across `N` nodes in the system to ensure high availability and fault tolerance.
- **Read/Write Requirements:** A successful `GET` operation requires responses from at least `R` nodes, while a `WRITE` operation needs at least `W` nodes to respond.
- **Dynamic Node Handling:** The system can handle the addition or removal of nodes dynamically, with automatic reassignment of key ranges and data migration.
- **Hinted Handoff:** If a node supposed to store a replica is down, the data is temporarily stored in a subsequent node, which continuously checks for the original node to come back online and transfers the data back to it once it does.
- **Gossip-based Protocol:** Nodes use a gossip-based protocol to propagate the membership list, ensuring system scalability and robustness.

# How to get it running
You can build the entire network using the `docker compose` file

It contains a `webclient` which is a simple web client to interact with the network and `node` container which is the actual network

Just simply run `docker-compose up --build` in the terminal and it should build the network for you

If you wish to add more nodes to the network, you can simply add more nodes in the `docker-compose.yml` file with similar fields and run `docker-compose up --build` again

When a network is running and you wish to simulate a new node to join the network. you can construct the a `make` command that is similar to the `run-second-node` command in the `Makefile`


# How It Works

## Read Requests

1. Upon receiving a `GET` request from the client, the system pings random machines within the distributed data store.
2. The machine that responds the fastest becomes the coordinator for that request.
3. Using consistent hashing, the coordinator calculates the nodes supposed to handle the key and forwards the request to them.
4. The nodes retrieve the data and respond back to the coordinator, which then returns the result to the client.

## Write Requests

1. When a client issues a `WRITE` request, the coordinator node updates the data item's vector clock, reflecting the new version of the data.
2. The coordinator stores the new data with its updated vector clock and forwards it to other nodes responsible for its replication.
3. Replicating nodes save the data and resolve any potential conflicts, typically using the last-write-wins principle.

## Dynamic Nodes Management

- **Adding Nodes:** When a new node joins, key ranges are reassigned based on the hash value of the new node's port number. Data stored in replicas that fall within the new node's key range are migrated to it.
- **Removing Nodes:** When a node is removed, its key ranges are released, and other nodes in the system take responsibility for them. Data in replicas stored in the removed node are transferred to the nodes that take over its key ranges.

## Fault Tolerance

- **Hinted Handoff:** If a node expected to store a replica is down, another node temporarily stores the data until the original node is back online, ensuring data availability during failures.

# Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Prerequisites

- setup [grpc](https://grpc.io/docs/languages/go/quickstart/) for go


## Installation

Provide step-by-step series of examples to help the user get the development environment running.

```sh
git clone [repo-url]
cd [repo-directory]
# additional commands...
```

## Configuration

Explain any configuration files/settings that need to be adapted to the local environment.

# Usage

Regenerate `proto` files, run `make proto`


# Contributing

Instructions on how to contribute to the project, submit pull requests, and report bugs.

# License

This project is licensed under the [LICENSE_NAME] - see the LICENSE file for details.
