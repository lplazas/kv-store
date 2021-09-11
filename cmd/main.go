package main

import (
	"context"
	"fmt"
	"github.com/gc-plazas/kv-store/internal"
)

func main() {
	c, err := internal.NewCluster(5, 10)
	ctx := context.TODO()
	if err := c.PutValue(ctx, "oslo", "norway"); err != nil {
		panic(err)
	}
	if err := c.PutValue(ctx, "cop", "den"); err != nil {
		panic(err)
	}
	if val, err := c.GetValue(ctx, "oslo"); err != nil {
		panic(err)
	} else {
		fmt.Println("val", val)
	}

	if val, err := c.GetValue(ctx, "osslo"); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		fmt.Println("val", val)
	}
	if err != nil {
		panic(err)
	}
}
