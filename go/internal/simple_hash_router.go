package internal

import (
	"fmt"
	"hash"
	"hash/fnv"
)

type simpleHashRouter struct {
	h hash.Hash32
}

func NewSimpleHashRouter() DocumentRouter {
	return simpleHashRouter{h: fnv.New32()}
}

func (r simpleHashRouter) GetShardNumber(documentID string, numShards int) (int, error) {
	r.h.Reset()
	_, err := r.h.Write([]byte(documentID))
	if err != nil {
		return 0, fmt.Errorf("fail generating hash, err: %w", err)
	}

	return int(r.h.Sum32()) % numShards, nil
}
