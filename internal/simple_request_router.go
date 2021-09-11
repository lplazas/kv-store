package internal

import (
	"context"
	"errors"
	"github.com/gc-plazas/kv-store/internal/errs"
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

	for !successful {
		randomPosition := rand.Intn(len(r.nodes))
		targetNode := r.nodes[randomPosition]
		nodeResult, err := targetNode.service.GetValue(ctx, shardID, key)
		successful = err == nil
		if !successful && errors.As(err, &errs.ValueNotFound{}) {
			return "", err
		} else {
			resultValue = nodeResult
		}
		r.removeNodeByIndex(randomPosition)
	}

	return resultValue, nil
}

func (r *SimpleRequestRouter) removeNodeByIndex(i int) {
	r.nodes[i] = r.nodes[len(r.nodes)-1]
	r.nodes = r.nodes[:len(r.nodes)-1]
}
