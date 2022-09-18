package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ids := fetchWebData(ctx)
	fmt.Println(ids)

}

func fetchWebData(ctx context.Context) (res []int64) {
	select {
	case <-time.After(3 * time.Second):
		return []int64{100, 200, 300}
	case <-ctx.Done():
		return []int64{1, 2, 3}
	}
}
