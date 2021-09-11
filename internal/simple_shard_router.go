package internal

import (
	"fmt"
	"github.com/gc-plazas/kv-store/internal/errs"
)

type SimpleShardRouter struct{}

func (r *SimpleShardRouter) RouteShardToNode(_ *Shard, nodes []*Node, excludeNodes []*Node) (*Node, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("cant route to an empty node list")
	}

	var eligibleNodes []*Node
AllNodeLoop:
	for _, node := range nodes {
		for _, excludeNode := range excludeNodes {
			if node.id == excludeNode.id {
				continue AllNodeLoop
			}
		}

		if node.service.IsHealthy() {
			eligibleNodes = append(eligibleNodes, node)
		}
	}

	if len(eligibleNodes) == 0 {
		return nil, errs.TryAgainLaterError("no healthy eligible nodes available")
	}

	leastShardCount := len(eligibleNodes[0].shards)
	nodeLeastCount := eligibleNodes[0]
	for _, n := range eligibleNodes {
		if len(n.shards) < leastShardCount {
			leastShardCount = len(n.shards)
			nodeLeastCount = n
		}
	}
	return nodeLeastCount, nil
}
