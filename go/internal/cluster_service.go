package internal

import (
	"context"
	"fmt"
	errs2 "github.com/gc-plazas/kv-store/go/internal/errs"
	"github.com/hashicorp/go-multierror"
	"strconv"
	"sync"
)

type RequestRouter interface {
	RouteRequest(nodesWithShard []*Node) *Node
}

func (c *Cluster) GetValue(ctx context.Context, key string) (string, error) {
	shardNum, err := c.docRouter.GetShardNumber(key, c.primaryShardCount)
	if err != nil {
		return "", errs2.FatalError("could not generate target shard id", err)
	}
	shardID := strconv.Itoa(shardNum)

	nodesWithShard := findNodesWithShard(c.nodes, shardID)

	if len(nodesWithShard) == 0 {
		// trigger cluster health check
		return "", errs2.TryAgainLaterError("data not available, try again later")
	}

	value, err := NewSimpleRequestRouter(nodesWithShard).RouteGetRequest(ctx, shardID, key)
	if err != nil {
		// trigger cluster health check
		return "", err
	}

	return value, nil
}

func (c *Cluster) PutValue(ctx context.Context, key, value string) error {
	shardNum, err := c.docRouter.GetShardNumber(key, c.primaryShardCount)
	if err != nil {
		return errs2.FatalError("could not generate shard number", err)
	}
	shardID := strconv.Itoa(shardNum)

	var replicaNodes []*Node
	var primaryNode *Node
	for _, node := range c.nodes {
		for _, allocatedShard := range node.shards {
			if allocatedShard.shard.id == shardID && !allocatedShard.shard.replica {
				primaryNode = node
			} else if allocatedShard.shard.id == shardID {
				replicaNodes = append(replicaNodes, node)
			}
		}
	}

	if primaryNode == nil {
		//trigger replica promotion
		//comeback later (?) wait for replica promotion?
		return errs2.FatalError("primary replica for value not found", nil)
	}

	if len(replicaNodes) == 0 {
		//trigger replica creation
	}

	if err := primaryNode.PutValue(ctx, shardID, key, value); err != nil {
		return fmt.Errorf("failed Put replicaNode, err: %w", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(replicaNodes))
	for _, replicaNode := range replicaNodes {
		n := replicaNode
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := n.PutValue(ctx, shardID, key, value); err != nil {
				errChan <- fmt.Errorf("failed Put replicaNode, err: %w", err)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled")
	case err, open := <-errChan:
		if !open {
			return nil
		}
		for _err := range errChan {
			err = multierror.Append(err, _err)
		}
		return fmt.Errorf("failed writing to replicas, err: %w", err)
	}
}
