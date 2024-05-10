package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const bufSize = 5

func main() {
	ch1 := make(chan int, bufSize)
	ch2 := make(chan int, bufSize)
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()
	wg.Add(3)
	go countProducer(&wg, ch1, bufSize, 50)
	go countProducer(&wg, ch2, bufSize, 500)
	go countConsumer(ctx, &wg, ch1, ch2)
	wg.Wait()
	fmt.Println("finish")
}

// chan<-int は書き込み専用channel型
func countProducer(wg *sync.WaitGroup, ch chan<- int, size int, sleep int) {
	defer wg.Done()
	defer close(ch)
	for i := 0; i < size; i++ {
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		ch <- i
	}
}

// <-chan int は読み込み専用channel型
func countConsumer(ctx context.Context, wg *sync.WaitGroup, ch1 <-chan int, ch2 <-chan int) {
	defer wg.Done()
	//loop:
	for ch1 != nil || ch2 != nil {
		select {
		case <-ctx.Done(): // タイムアウトが発生した場合
			fmt.Println(ctx.Err())
			//break loop
			// タイムアウト時にチャネルに残っている値をすべて読み込むように変更
			for ch1 != nil || ch2 != nil {
				select {
				case v, ok := <-ch1: //ch1に値が入っている場合または、closeされている場合
					if !ok { //close
						ch1 = nil
						break // caseを抜ける
					}
					fmt.Printf("ch1 %v\n", v)
				case v, ok := <-ch2:
					if !ok {
						ch2 = nil
						break
					}
					fmt.Printf("ch2 %v\n", v)
				}
			}
		case v, ok := <-ch1: //ch1に値が入っている場合または、closeされている場合
			if !ok { //close
				ch1 = nil
				break // caseを抜ける
			}
			fmt.Printf("ch1 %v\n", v)
		case v, ok := <-ch2:
			if !ok {
				ch2 = nil
				break
			}
			fmt.Printf("ch2 %v\n", v)
		}
	}
}
