package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(<-ch1)
	}()
	ch1 <- 10
	close(ch1)
	v, ok := <-ch1
	fmt.Printf("%v %v\n", v, ok) // 0 false
	wg.Wait()

	ch2 := make(chan int, 2)
	ch2 <- 1
	ch2 <- 2
	close(ch2)
	v, ok = <-ch2
	fmt.Printf("%v %v\n", v, ok) // 1 true
	v, ok = <-ch2
	fmt.Printf("%v %v\n", v, ok) // 2 true
	v, ok = <-ch2
	fmt.Printf("%v %v\n", v, ok) // 0 false

	// カプセル化の例
	ch3 := generateCountStream()
	for v := range ch3 {
		fmt.Println(v)
	}

	// 通知専用のchannelの例
	nCh := make(chan struct{}) // 空の構造体のサイズは0byte
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("goroutine %v started\n", i)
			<-nCh
			fmt.Println(i)
		}(i)
	}
	time.Sleep(2 * time.Second)
	close(nCh)
	fmt.Println("unlocked by manual close")
	wg.Wait()
	fmt.Println("finish")
}

// <-chanは読み取り専用のchannel型
func generateCountStream() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i <= 5; i++ {
			ch <- i
		}
	}()
	return ch
}
