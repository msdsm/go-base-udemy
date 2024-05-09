package main

import (
	"fmt"
	"runtime"
)

func main() {
	// ch := make(chan int)
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	ch <- 10
	// 	time.Sleep(500 * time.Millisecond)
	// }()
	// fmt.Println(<-ch)
	// wg.Wait()

	ch1 := make(chan int)
	go func() { // goroutine leak
		fmt.Println(<-ch1)
	}()
	ch1 <- 10
	fmt.Printf("num of working goroutines: %d\n", runtime.NumGoroutine())

	ch2 := make(chan int, 1)
	ch2 <- 2
	ch2 <- 2
	fmt.Println(<-ch2)
}
