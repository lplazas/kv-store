# Key considerations for KV store

## HA - Scalability
- Sharding -> without knowledge of the kind of key it is better to shard by hash of key -> consistent hashing
- Replication 
  - Need to implement primary/replica pattern for shards
  - Index always to primary, replicate to replicas
  - Mark as ready when all have applied a change
  - If primary shard is down, promote a replica
- Need of a leader to keep cluster state
- Routing/allocation algorithms
## Durability
- Different layers/tiers of storage (cached/disk)
## Reliability
- Error handling 
  - Handle node failure
  - Handle shard failure
  - Handle write failure
  - Route to healthy node/shard
## Consistency
- Transactions ?
- Entry versioning ?

# Goals for this simple implementation :

- Represent nodes as separate in memory structs
- Single table approach
- In memory map used as backing store 
- Main goroutine to act as master and track cluster state
- Only main node to receive requests
- Channels as node-node communication 
- Constant number of nodes=5 and shards=10 (1- primary and 10 replicas)
- Handle some errors 
  - Node unavailable ->
    - reroute request to replica
    - turn replica into primary
    - Create a new replica in another node
  
# How to run 

1. Clone this repo
2. Install docker: https://docs.docker.com/get-docker/
3. Build the binaries by running: `make build-server && make build-client`
4. Run the server with: `docker run -p 1338:1338 kv-store-server:latest`
5. Run the client with: 
```
docker run --network="host" kv-store-client:latest localhost:1338 put name luis
docker run --network="host" kv-store-client:latest localhost:1338 get name
```
