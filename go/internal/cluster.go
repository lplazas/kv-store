package internal

import (
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/gc-plazas/kv-store/go/internal/errs"
	"strconv"
)

type DocumentRouter interface {
	// GetShardNumber - Routes docIDs uniformly to a shard between 0 and numShards-1
	GetShardNumber(documentID string, numShards int) (int, error)
}

type ShardRouter interface {
	// RouteShardToNode - Routes shards to a node with the least # of shards
	RouteShardToNode(shard *Shard, nodes []*Node, excludeNodes []*Node) (*Node, error)
}

type Cluster struct {
	nodeCount         int
	primaryShardCount int
	nodes             []*Node
	shards            []*Shard
	docRouter         DocumentRouter
	shardRouter       ShardRouter
	requestRouter     RequestRouter
}

func NewCluster(nodeCount, primaryShardCount int) (*Cluster, error) {
	cluster := &Cluster{
		nodeCount:         nodeCount,
		primaryShardCount: primaryShardCount,
		docRouter:         NewSimpleHashRouter(),
		shardRouter:       &SimpleShardRouter{},
	}

	cluster.initializeNodes()

	if err := cluster.initializePrimaryShards(); err != nil {
		return nil, fmt.Errorf("failed shard initialization, err: %w", err)
	}

	if err := cluster.initializeReplicaShards(); err != nil {
		return nil, fmt.Errorf("failed shard initialization, err: %w", err)
	}

	cluster.printClusterState()

	return cluster, nil
}

func findNodesWithShard(nodes []*Node, shardID string) []*Node {
	var nodesWithShard []*Node
	for _, node := range nodes {
		for _, allocatedShard := range node.shards {
			if allocatedShard.shard.id == shardID {
				nodesWithShard = append(nodesWithShard, node)
			}
		}
	}
	return nodesWithShard
}

func (c *Cluster) initializePrimaryShards() error {
	for i := 0; i < c.primaryShardCount; i++ {
		shard := NewShard(strconv.Itoa(i), false)
		c.shards = append(c.shards, shard)

		if err := backoff.Retry(func() error {
			nodeToAllocate, err := c.shardRouter.RouteShardToNode(shard, c.nodes, nil)

			if err != nil {
				return translateBackoffErr(err)
			}

			err = nodeToAllocate.AllocateShard(shard, map[string]string{})

			if err != nil {
				return translateBackoffErr(err)
			}

			shard.AllocateToNode(nodeToAllocate.id)

			return nil
		}, backoff.NewExponentialBackOff()); err != nil {
			return fmt.Errorf("failed primary shard initialization, err: %w", err)
		}
	}
	return nil
}

func (c *Cluster) initializeReplicaShards() error {
	for _, primaryShard := range c.shards {
		replica := NewShard(primaryShard.id, true)
		c.shards = append(c.shards, replica)

		if err := backoff.Retry(func() error {
			excludeAllocationNodes := []*Node{c.getNode(primaryShard.allocatedNodeID)}
			nodeToAllocate, err := c.shardRouter.RouteShardToNode(replica, c.nodes, excludeAllocationNodes)

			if err != nil && errors.As(err, &errs.TryAgainLater{}) {
				return err
			} else if err != nil {
				return &backoff.PermanentError{Err: err}
			}

			err = nodeToAllocate.AllocateShard(replica, map[string]string{})

			if err != nil && errors.As(err, &errs.TryAgainLater{}) {
				return err
			} else if err != nil {
				return &backoff.PermanentError{Err: err}
			}

			replica.AllocateToNode(nodeToAllocate.id)
			return nil
		}, backoff.NewExponentialBackOff()); err != nil {
			return fmt.Errorf("failed primary primaryShard initialization, err: %w", err)
		}
	}
	return nil
}

func (c *Cluster) printClusterState() {
	for _, node := range c.nodes {
		fmt.Println("node -", node.id)
		for _, shard := range node.shards {
			fmt.Println("   shard -", shard.shard.id, "replica:", shard.shard.replica)
		}
	}
}

//normally nodes would -join- the cluster
func (c *Cluster) initializeNodes() {
	for i := 0; i < c.nodeCount; i++ {
		n := NewNode(strconv.Itoa(i))
		c.nodes = append(c.nodes, n)
	}
}

func translateBackoffErr(err error) error {
	if err != nil && !errors.As(err, &errs.TryAgainLater{}) {
		return &backoff.PermanentError{Err: err}
	}
	return err
}

func (c *Cluster) getNode(nodeID string) *Node {
	for _, node := range c.nodes {
		if node.id == nodeID {
			return node
		}
	}
	return nil
}
