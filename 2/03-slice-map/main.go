package main

import "fmt"

func main() {
	// 配列
	var a1 [3]int // すべて0で初期化
	var a2 = [3]int{10, 20, 30}
	a3 := [...]int{10, 20} // [2]int型であることを推論してくれる
	fmt.Printf("%v %v %v\n", a1, a2, a3)
	fmt.Printf("%v %v\n", len(a3), cap(a3))
	fmt.Printf("%T %T\n", a2, a3)

	// スライス
	var s1 []int
	s2 := []int{}
	fmt.Printf("s1: %[1]T %[1]v %v %v\n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %[1]T %[1]v %v %v\n", s2, len(s2), cap(s2))
	fmt.Println(s1 == nil)   // これはtrue
	fmt.Println(s2 == nil)   // これはfalse
	s1 = append(s1, 1, 2, 3) // 要素数の追加方法
	fmt.Printf("s1: %[1]T %[1]v %v %v\n", s1, len(s1), cap(s1))
	s3 := []int{4, 5, 6}
	s1 = append(s1, s3...) // スライスの要素をすべて別のスライスの末尾に追加する方法
	fmt.Printf("s1: %[1]T %[1]v %v %v\n", s1, len(s1), cap(s1))

	s4 := make([]int, 0, 2) // スライスの型, 要素数, キャパシティを宣言
	fmt.Printf("s4: %[1]T %[1]v %v %v\n", s4, len(s4), cap(s4))
	s4 = append(s4, 1, 2, 3, 4)
	fmt.Printf("s4: %[1]T %[1]v %v %v\n", s4, len(s4), cap(s4))

	s5 := make([]int, 4, 6)
	fmt.Printf("s5: %v %v %v\n", s5, len(s5), cap(s5))
	s6 := s5[1:3] // [l, r)の半開区間
	s6[1] = 10    // スライスは参照型なのでs5[2]も10になる
	fmt.Printf("s5: %v %v %v\n", s5, len(s5), cap(s5))
	fmt.Printf("s6: %v %v %v\n", s6, len(s6), cap(s6))
	s6 = append(s6, 2)
	fmt.Printf("s5: %v %v %v\n", s5, len(s5), cap(s5))
	fmt.Printf("s6 appended: %v %v %v\n", s6, len(s6), cap(s6))

	// 上のようにコピー元の値を変更させない方法
	// 新しくスライスを作成(新しいメモリ領域を確保してそこを参照先とする)
	sc6 := make([]int, len(s5[1:3]))
	fmt.Printf("s5 source of copy: %v %v %v\n", s5, len(s5), cap(s5))
	fmt.Printf("sc6 dst copy before: %v %v %v\n", sc6, len(sc6), cap(sc6))
	copy(sc6, s5[1:3]) // 参照のコピーではなく値のコピーになる copy(src,dst)に対してdstの該当する範囲の値をsrcの参照先に格納するということ
	fmt.Printf("sc6 dst of copy after: %v %v %v\n", sc6, len(sc6), cap(sc6))
	sc6[1] = 12 // makeしてcopyしたのでs5は不変
	fmt.Printf("s5: %v %v %v\n", s5, len(s5), cap(s5))
	fmt.Printf("sc6: %v %v %v\n", sc6, len(sc6), cap(sc6))

	// メモリの共有を部分的に許可する方法
	s5 = make([]int, 4, 6)
	// スライス[l:r:x]でメモリの共有を部分的に許可できる
	// コピー元の[l,r)区間の値を参照するが、コピー元の[0, x)のみメモリを共有する
	fs6 := s5[1:3:3]
	fmt.Printf("s5: %v %v %v\n", s5, len(s5), cap(s5))
	fmt.Printf("fs6: %v %v %v\n", fs6, len(fs6), cap(fs6))
	fs6[0] = 6
	fs6[1] = 7
	fs6 = append(fs6, 8)
	fmt.Printf("s5: %v %v %v\n", s5, len(s5), cap(s5))
	fmt.Printf("fs6: %v %v %v\n", fs6, len(fs6), cap(fs6))
	s5[3] = 9
	fmt.Printf("s5: %v %v %v\n", s5, len(s5), cap(s5))
	fmt.Printf("fs6: %v %v %v\n", fs6, len(fs6), cap(fs6))

	// map使い方
	// map[keyの型]valueの型
	var m1 map[string]int
	m2 := map[string]int{}
	fmt.Printf("%v %v \n", m1, m1 == nil) // これはtrue
	fmt.Printf("%v %v \n", m2, m2 == nil) // これはfalse
	m2["A"] = 10
	m2["B"] = 20
	m2["C"] = 0
	// lenで(key,value)の個数を取得
	fmt.Printf("%v %v %v\n", m2, len(m2), m2["A"])
	// keyを指定して(key,value)を削除
	delete(m2, "A")
	fmt.Printf("%v %v %v\n", m2, len(m2), m2["A"])
	v, ok := m2["A"] // 存在しないなら返り値の2つ目がfalse
	fmt.Printf("%v %v\n", v, ok)
	v, ok = m2["C"]
	fmt.Printf("%v %v\n", v, ok)

	// 全要素取得
	// mapは内部でハッシュを使っているから順番が同じになるとは限らない
	for k, v := range m2 {
		fmt.Printf("%v %v\n", k, v)
	}
}
