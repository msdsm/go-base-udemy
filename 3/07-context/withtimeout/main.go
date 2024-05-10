// WithTimeoutの例
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	wg.Add(3)
	go subTask(ctx, &wg, "a")
	go subTask(ctx, &wg, "b")
	go subTask(ctx, &wg, "c")
	wg.Wait()
}

func subTask(ctx context.Context, wg *sync.WaitGroup, id string) {
	defer wg.Done()
	// 一定間隔でチャネルの書き込みを発生させる
	t := time.NewTicker(500 * time.Millisecond)
	select {
	case <-ctx.Done(): // タイムアウト発生
		fmt.Println(ctx.Err())
		return
	case <-t.C: //よみこみ
		t.Stop() //定時生成をとめる
		fmt.Println(id)
	}
}
