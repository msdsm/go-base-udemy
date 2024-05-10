package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)
		ch1 <- "A"
	}()
	go func() {
		defer wg.Done()
		time.Sleep(800 * time.Millisecond)
		ch2 <- "B"
	}()
	cnt := 0
loop:
	for ch1 != nil || ch2 != nil {
		cnt++
		select {
		case <-ctx.Done():
			fmt.Println("Timeout")
			break loop
		case v := <-ch1:
			fmt.Println(v)
			ch1 = nil
		case v := <-ch2:
			fmt.Println(v)
			ch2 = nil
		}
	}
	wg.Wait()
	fmt.Println("finish")
	fmt.Println(cnt)
}
