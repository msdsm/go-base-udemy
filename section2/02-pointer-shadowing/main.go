package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// アドレス取得
	var ui1 uint16
	fmt.Printf("memory address of uil: %p\n", &ui1)
	var ui2 uint16
	fmt.Printf("memory address of uil: %p\n", &ui2)

	// ポインタ変数宣言
	var p1 *uint16 // nilで初期化される
	fmt.Printf("value of p1: %v\n", p1)
	p1 = &ui1 // 代入
	fmt.Printf("value of p1: %v\n", p1)
	// ポインタ変数が格納されている番地
	fmt.Printf("size of p1: %d[bytes]\n", unsafe.Sizeof(p1))
	fmt.Printf("memory address of p1: %p\n", &p1)

	// dereference
	fmt.Printf("value of ui1(dereference): %v\n", *p1)
	*p1 = 1
	fmt.Println(ui1)

	// ポインタのポインタ
	var pp1 **uint16 = &p1
	fmt.Printf("value of pp1: %p\n", &pp1)
	fmt.Printf("size of pp1: %d[bytes]\n", unsafe.Sizeof(pp1))
	fmt.Printf("value of p1(dereference): %v\n", *pp1)   // 1回dereference
	fmt.Printf("value of ui1(dereference): %v\n", **pp1) // 2回dereference

	// シャドーイングの確認
	ok, res := true, "A"
	fmt.Printf("memory address of res: %p\n", &res) // res変数のメモリ確認
	if ok {
		// シャドーイング
		// :=を使うことでローカル変数を定義、このブロック外で使用されている名前と衝突しない
		res := "B"
		fmt.Printf("memory address of res: %p\n", &res) // メモリ確認(ブロック外のresと違うことがわかる)
	}
	// あたりまえだけど以下の場合は同じメモリ
	if ok {
		res = "C"
		fmt.Printf("memory address of res: %p\n", &res) // メモリ確認(ブロック外のresと違うことがわかる)
	}
}
