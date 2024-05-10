package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Millisecond))
	defer cancel()
	ch := subTask(ctx)
	v, ok := <-ch
	if ok {
		fmt.Println(v)
	}
	fmt.Println("finish")
}

func subTask(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		deadline, ok := ctx.Deadline() // deadlineの時刻取得
		if ok {                        // deadlineが設定されているなら
			if deadline.Sub(time.Now().Add(30*time.Millisecond)) < 0 { // deadlineに次に実行するタスクが間に合うかどうか
				fmt.Println("impossible to meet deadline") // 間に合わない
				return                                     // ここで終了するまでにdeferでchannelがcloseされる
			}
		}
		// 間に合うまたはdeadlineが設定されていないならsubtask実行
		time.Sleep(30 * time.Millisecond)
		ch <- "hello"
	}()
	return ch
}
