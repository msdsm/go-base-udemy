package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		ch1 := make(chan string, 1)
		ch2 := make(chan string, 1)
		ch1 <- "ch1"
		ch2 <- "ch2"
		select {
		case v := <-ch1:
			fmt.Println(v)
		case v := <-ch2:
			fmt.Println(v)
		}
	}
}
