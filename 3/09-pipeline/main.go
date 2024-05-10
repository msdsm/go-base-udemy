/*
実行結果
8
12
16
finish
*/

package main

import (
	"context"
	"fmt"
)

func generator(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case <-ctx.Done(): // cancel用
				return
			case out <- n:
			}
		}
	}()
	return out
}

func double(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in { // channelに書き込みがあるたびに値を読み込む
			select {
			case <-ctx.Done(): // cancel用
				return
			case out <- n * 2:
			}
		}
	}()
	return out
}

func offset(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in { // channelに書き込みがあるたびに値を読み込む
			select {
			case <-ctx.Done(): // cancel用
				return
			case out <- n + 2:
			}
		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nums := []int{1, 2, 3, 4, 5}
	var i int
	flag := true

	// 値が書き込まれるたびに実行(channelがcloseされると抜ける)
	for v := range double(ctx, offset(ctx, double(ctx, generator(ctx, nums...)))) {
		if i == 3 {
			cancel()
			flag = false
		}
		if flag {
			fmt.Println(v)
		}
		i++
	}
	fmt.Println("finish")
}
