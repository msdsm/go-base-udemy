package main

import (
	"fmt"
	"time"
)

func main() {
	a := -1
	if a == 0 {
		fmt.Println("zero")
	} else if a > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("negative")
	}
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// for{}で無限ループ
	// for {
	// 	fmt.Println("working")
	// 	time.Sleep(2 * time.Second)
	// }

	// 無限ループを抜けるにはbreak
	var i int
	for {
		if i > 3 {
			break
		}
		fmt.Println(i)
		i += 1
		time.Sleep(300 * time.Millisecond)
	}

	// breakで抜けるブロックを明示する方法
loop: // forブロックにラベルつける
	for i := 0; i < 10; i++ {
		switch i {
		case 2:
			continue
		case 3:
			continue
		case 8:
			break loop // これでloopラベルのついたfor文から抜ける(こうしないとswitch文を抜けるだけになる)
		default:
			fmt.Printf("%v ", i)
		}
	}
	fmt.Printf("\n")

	// 構造体のスライス定義
	items := []item{
		{price: 10.},
		{price: 20.},
		{price: 30.},
	}
	// (index, element)をrangeで取得
	for _, item := range items {
		item.price *= 1.1 // この(index,element)は構造体の値をコピーであって参照ではないので値変わらない
	}
	fmt.Printf("%+v\n", items) // priceは変わっていない
	for i := range items {
		items[i].price *= 1.1 // 値を変えたい場合はこう
	}
	fmt.Printf("%+v\n", items)
}

type item struct {
	price float32
}
