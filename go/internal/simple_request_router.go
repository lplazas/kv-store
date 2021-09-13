package internal

import (
	"context"
	"errors"
	"github.com/gc-plazas/kv-store/go/internal/errs"
	"math/rand"
)

type SimpleRequestRouter struct {
	nodes []*Node
}

func NewSimpleRequestRouter(nodes []*Node) *SimpleRequestRouter {
	return &SimpleRequestRouter{nodes: nodes}
}

func (r *SimpleRequestRouter) RouteGetRequest(ctx context.Context, shardID, key string) (string, error) {
	var resultValue string
	successful := false

	var eligibleNodes []*Node
	for _, node := range r.nodes {
		if node.IsHealthy() {
			eligibleNodes = append(eligibleNodes, node)
		}
	}

	for !successful {
		randomPosition := rand.Intn(len(eligibleNodes))
		targetNode := eligibleNodes[randomPosition]
		nodeResult, err := targetNode.GetValue(ctx, shardID, key)
		successful = err == nil
		if !successful && errors.As(err, &errs.ValueNotFound{}) {
			return "", err
		} else {
			resultValue = nodeResult
		}
		removeNodeByIndex(eligibleNodes, randomPosition)
	}

	return resultValue, nil
}

func removeNodeByIndex(nodes []*Node, i int) {
	nodes[i] = nodes[len(nodes)-1]
	nodes = nodes[:len(nodes)-1]
}
