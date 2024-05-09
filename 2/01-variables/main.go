package main

import "fmt"

const secret = "abc"

type Os int

// const宣言のiotaを用いた連番の使用例
const (
	Mac Os = iota + 1
	Windows
	Linux
)

// 変数の一括宣言
var (
	in int
	st string
	bl bool
)

func main() {
	// var i int
	// var i int = 2
	// var i = 2
	i := 1
	ui := uint16(2)
	fmt.Println(i)
	fmt.Printf("i: %v %T\n", i, i)
	fmt.Printf("i: %[1]v %[1]T ui: %[2]v %[2]T\n", i, ui)

	f := 1.23456
	s := "hello"
	b := true
	fmt.Printf("f: %[1]v %[1]T\n", f)
	fmt.Printf("s: %[1]v %[1]T\n", s)
	fmt.Printf("b: %[1]v %[1]T\n", b)

	// 同時宣言
	pi, title := 3.14, "Go"
	fmt.Printf("pi: %v title: %v\n", pi, title)

	// 型変換
	x := 10
	y := 1.23
	// z := x + y // これはできない
	z := float64(x) + y
	fmt.Println(z)

	fmt.Printf("Mac:%v Wdinwos:%v Linux:%v\n", Mac, Windows, Linux)
}
