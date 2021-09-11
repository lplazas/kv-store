package internal

import (
	"hash"
	"strconv"
	"testing"
)

func TestSimpleHashRouter_GetShardNumber(t *testing.T) {
	type fields struct {
		h hash.Hash32
	}
	type args struct {
		documentID string
		numShards  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "3 shards",
			args: args{
				documentID: "testDocId",
				numShards:  3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSimpleHashRouter()
			distributionMap := make(map[int]int)
			for i := 0; i < 1000; i++ {
				got, err := r.GetShardNumber(strconv.Itoa(i)+tt.args.documentID, tt.args.numShards)
				distributionMap[got]++
				if err != nil {
					panic(err)
				}
			}

			if len(distributionMap) != tt.args.numShards {
				t.Errorf("Not generating entries between 0 and %d, got = %d, want %d", tt.args.numShards, len(distributionMap), tt.args.numShards)
			}
		})
	}
}
