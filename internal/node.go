package internal

import (
	"context"
	"github.com/gc-plazas/kv-store/internal/errs"
	"github.com/gc-plazas/kv-store/internal/inmem"
)

type NodeState int

const (
	NodeStateHealthy   = 0
	NodeStateUnhealthy = 1
)

type SimpleService struct {
	node *Node
}

type NodeService interface {
	GetValue(ctx context.Context, shardID, key string) (string, error)
	PutValue(ctx context.Context, shardID, key, value string) error
	MakeShardPrimary(shardID string) error
	GetAllFromShard(shardID string) ([]string, error)
	IsHealthy() bool
	AllocateShard(shard *Shard, values map[string]string) error
}

type AllocatedShard struct {
	shard   *Shard
	storage KeyValueStorage
}

type Node struct {
	id      string
	state   NodeState
	shards  []*AllocatedShard
	service NodeService
}

func NewNode(id string) *Node {
	node := &Node{
		id:    id,
		state: NodeStateHealthy,
	}
	node.service = SimpleService{node: node} //todo this is ugly
	return node
}

func (s SimpleService) GetValue(_ context.Context, shardID, key string) (string, error) {
	targetShard := s.findAllocatedShard(shardID)

	if targetShard == nil {
		return "", errs.NodeNotAvailable{}
	}

	value, err := targetShard.storage.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (s SimpleService) PutValue(_ context.Context, shardID, key, value string) error {
	targetShard := s.findAllocatedShard(shardID)

	if targetShard == nil {
		return errs.NodeNotAvailable{}
	}

	if err := targetShard.storage.Put(key, value); err != nil {
		return err
	}

	return nil
}

func (s SimpleService) findAllocatedShard(shardID string) *AllocatedShard {
	var targetShard *AllocatedShard
	for _, allocatedShard := range s.node.shards {
		if allocatedShard.shard.id == shardID {
			targetShard = allocatedShard
		}
	}
	return targetShard
}

func (s SimpleService) MakeShardPrimary(shardID string) error {
	panic("implement me")
}

func (s SimpleService) GetAllFromShard(shardID string) ([]string, error) {
	panic("implement me")
}

func (s SimpleService) IsHealthy() bool {
	return true
}

func (s SimpleService) AllocateShard(shard *Shard, values map[string]string) error {
	allocatedShard := &AllocatedShard{
		shard:   shard,
		storage: inmem.NewMemoryKeyValueStore(),
	}

	for k, v := range values {
		if err := allocatedShard.storage.Put(k, v); err != nil {
			return errs.TryAgainLaterError("failed allocating shard")
		}
	}

	s.node.shards = append(s.node.shards, allocatedShard)

	return nil
}
