# Authenticated Keyrowd Search on Large-scale Graphs in Hybrid-Storage Blockchains

## Introduction

Keyword search on graphs, which can be informally described as follows: Given a graph G and a query Q associated with a set of keywords, the keyword search aims to find rooted trees RT in G such that nodes in these trees collectively cover part of or all keywords in Q.
Keyword search is an important operation for analyzing graph data, providing a simple and user-friendly interface for retrieving information from complex graph data structures without prior knowledge of a specialized query language. The data owner could delegate their data graphs to the Tamper-proof blockchain due to the limited storage and computation. In the blockchain, to improve scalability, it is a good choice to store the raw data in an off-chain storage service provider (SP) and only maintain the digest of the raw data on-chain by smart contracts. However, if the SP is untrusted, it may return tampered-with results. To support integrity-assured data queries in such a scenario, authenticated subgraph matching can be applied. To our knowledge, there is no work to enable the blockchain to support keyword searches on graphs.

In this paper, we study a novel approach to support authenticated keyword search for the large graph kept off-chain. We first design a two-stage baseline solution to solve the problem. To reduce the evaluation time and VO size, a novel ADS MP-tree is proposed to handle authenticated keyword searches on graphs effectively. The main idea of MP-tree is to aggregate illegal paths that do not appear in the result trees. This comes with the observation that the result trees are composed of the shortest paths from the root node to the leaf nodes, i.e., a result tree consists of paths starting from the same root node and arriving at the leaf nodes that cover all the query keywords. Compared with the naive method, MP-tree can realize authenticated keyword searches on graphs effectively, but the index size of MP-tree is still considerable. To reduce the storage cost of the ADS, we proposed an optimized ADS MP-tree*, which combines similar subtrees in different MP-trees.
## Environment

Ethereum blockchain platform, Golang 1.15.5.

1. devDependencies

   ```
   CentOS:
   yum install git wget bzip2 vim gcc-c++ ntp epel-release
   nodejs cmake -y
   yum update
   Ubuntu:
   sudo apt install make
   sudo apt install g++
   sudo apt-get install libltdl-dev
   apt install -y build-essential
   ```

2. go-ethereum install

   ```
   git clone https://github.com/ethereum/go-ethereum
   cd go-ethereum
   make geth
   # Copy the finished go-ethereum /buil/bin/geth executable to /usr/local/bin
   ```

3. Create Genesis Block

    Create a new genesis.json file and write the following information

   ```
   {
   "config": {
   "chainId": 981106,
   "homesteadBlock": 0,
   "eip150Block": 0,
   "eip150Hash":
   "0x0000000000000000000000000000000000000000000000000000000000
   000000",
   "eip155Block": 0,
   "eip158Block": 0,
   "byzantiumBlock": 0,
   "constantinopleBlock": 0,
   "petersburgBlock": 0,
   "ethash": {}
   
   },
   "nonce": "0x0",
   "timestamp": "0x284d29c0",
   "extraData":
   "0x0000000000000000000000000000000000000000000000000000000000
   000000",
   "gasLimit": "0x47b760",
   "difficulty": "0x80000",
   "mixHash":
   "0x0000000000000000000000000000000000000000000000000000000000
   000000",
   "coinbase": "0x0000000000000000000000000000000000000000",
   "alloc": {
   "0000000000000000000000000000000000000000": {
   "balance": "0x1"
   }
   },
   "number": "0x0",
   "gasUsed": "0x0",
   "parentHash":
   "0x0000000000000000000000000000000000000000000000000000000000
   000000"
   }
   ```

## Test

1. Enabling the Network

   ```
   geth --datadir data init genesis.json
   geth --datadir data --networkid 981106 --http --
   http.corsdomain "*" --http.port 8545 --http.addr 0.0.0.0 --
   nodiscover console --allow-insecure-unlock
   ```

2. Open a new shell to mine

   ```
   geth attach ipc:geth.ipc
   miner.start()
   ```

3. Open a new shell to test

   ```
   cd src
   go run main.go
   ```

4. Remember to turn off the network after use, press ctrl+d or type 'exit'.

