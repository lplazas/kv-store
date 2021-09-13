package internal

import (
	"context"
	errs2 "github.com/gc-plazas/kv-store/go/internal/errs"
	"github.com/gc-plazas/kv-store/go/internal/inmem"
)

type NodeService interface {
	GetValue(ctx context.Context, shardID, key string) (string, error)
	PutValue(ctx context.Context, shardID, key, value string) error
	MakeShardPrimary(shardID string) error
	GetAllFromShard(shardID string) ([]string, error)
	IsHealthy() bool
	ChangeHealth(healthy bool)
	AllocateShard(shard *Shard, values map[string]string) error
}

func (n *Node) GetValue(_ context.Context, shardID, key string) (string, error) {
	targetShard := n.findAllocatedShard(shardID)

	if targetShard == nil {
		return "", errs2.NodeNotAvailable{}
	}

	value, err := targetShard.storage.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (n *Node) PutValue(_ context.Context, shardID, key, value string) error {
	targetShard := n.findAllocatedShard(shardID)

	if targetShard == nil {
		return errs2.NodeNotAvailable{}
	}

	if err := targetShard.storage.Put(key, value); err != nil {
		return err
	}

	return nil
}

func (n *Node) MakeShardPrimary(_ string) error {
	panic("implement me")
}

func (n *Node) GetAllFromShard(_ string) ([]string, error) {
	panic("implement me")
}

func (n *Node) IsHealthy() bool {
	return true
}

func (n *Node) ChangeHealth(healthy bool) {
	if healthy {
		n.state = NodeStateHealthy
	} else {
		n.state = NodeStateUnhealthy
	}
}

func (n *Node) AllocateShard(shard *Shard, values map[string]string) error {
	allocatedShard := &AllocatedShard{
		shard:   shard,
		storage: inmem.NewMemoryKeyValueStore(),
	}

	for k, v := range values {
		if err := allocatedShard.storage.Put(k, v); err != nil {
			return errs2.TryAgainLaterError("failed allocating shard")
		}
	}

	n.shards = append(n.shards, allocatedShard)

	return nil
}
