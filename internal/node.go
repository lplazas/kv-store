package internal

type NodeState int

const (
	NodeStateHealthy   = 0
	NodeStateUnhealthy = 1
)

type SimpleService struct {
	node *Node
}

type AllocatedShard struct {
	shard   *Shard
	storage KeyValueStorage
}

type Node struct {
	id     string
	state  NodeState
	shards []*AllocatedShard
}

func NewNode(id string) *Node {
	node := &Node{
		id:    id,
		state: NodeStateHealthy,
	}
	return node
}

func (n *Node) findAllocatedShard(shardID string) *AllocatedShard {
	var targetShard *AllocatedShard
	for _, allocatedShard := range n.shards {
		if allocatedShard.shard.id == shardID {
			targetShard = allocatedShard
		}
	}
	return targetShard
}
