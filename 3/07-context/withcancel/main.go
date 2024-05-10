package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		v, err := criticalTask(ctx)
		if err != nil {
			fmt.Printf("critical task cancelled due to: %v\n", err)
			cancel() // manual cancel
			return
		}
		fmt.Println("success", v)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		v, err := normalTask(ctx)
		if err != nil {
			fmt.Printf("normal task cancelled due to: %v\n", err)
			return
		}
		fmt.Println("success", v)
	}()
	wg.Wait()
}

func criticalTask(ctx context.Context) (string, error) {
	// 親のコンテキストをもとに新しくコンテキストを生成
	ctx, cancel := context.WithTimeout(ctx, 800*time.Millisecond)
	defer cancel() // go leakを防ぐため
	// 定時生成
	t := time.NewTicker(1000 * time.Millisecond)
	select {
	case <-ctx.Done(): // タイムアウト
		return "", ctx.Err()
	case <-t.C: // channel読み込み
		t.Stop() // 定時生成stop
	}

	// タイムアウトが発生しなかったときの返り値
	return "A", nil
}

func normalTask(ctx context.Context) (string, error) {
	t := time.NewTicker(3000 * time.Millisecond)
	select {
	case <-ctx.Done(): // これは親のcontext, つまりWithCancel
		return "", ctx.Err()
	case <-t.C:
		t.Stop()
	}
	return "B", nil
}
