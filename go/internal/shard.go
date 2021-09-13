package internal

type Shard struct {
	id              string
	replica         bool
	allocatedNodeID string
}

type KeyValueStorage interface {
	Get(key string) (string, error)
	Put(key, value string) error
}

func NewShard(id string, replica bool) *Shard {
	return &Shard{
		id:      id,
		replica: replica,
	}
}

func (s *Shard) AllocateToNode(nodeID string) {
	s.allocatedNodeID = nodeID
}
