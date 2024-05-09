package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// deferの動作確認
func funcDefer() {
	defer fmt.Println("main func final-finish")
	defer fmt.Println("main func semi-finish")
	fmt.Println("Hello world")
}

// 可変引数の使い方
// 複数引数で受け取ったものをfilesというスライスで扱える
func trimExtension(files ...string) []string {
	out := make([]string, 0, len(files))
	for _, f := range files { // index, element
		out = append(out, strings.TrimSuffix(f, ".csv")) // .csvを切る
	}
	return out
}

func fileChecker(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", errors.New("file not found")
	}
	defer f.Close()
	return name, nil
}

// 引数で無名関数を利用
func addExt(f func(file string) string, name string) {
	fmt.Println(f(name))
}

// 返り値が無名関数
func multiply() func(int) int {
	return func(n int) int {
		return n * 1000
	}
}

// closure
func countUp() func(int) int {
	count := 0
	return func(n int) int {
		count += n
		return count
	}
}

func main() {
	funcDefer()

	// スライス用意
	files := []string{"file1.csv", "file2.csv", "file3.csv"}
	// スライス変数名...でその要素を順に渡せる
	fmt.Println(trimExtension(files...))
	name, err := fileChecker("file.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(name)

	i := 1
	// 無名関数宣言、即時実行
	func(i int) {
		fmt.Println(i)
	}(i)

	// 無名関数を格納した変数
	f1 := func(i int) int {
		return i + 1
	}
	fmt.Println(f1(i))

	// 無名関数の変数宣言
	f2 := func(file string) string {
		return file + ".csv"
	}
	// 引数で無名関数を利用する例
	addExt(f2, "file1")

	// 返り値で無名関数を利用する例
	f3 := multiply()
	fmt.Println(f3(2))

	// closure
	f4 := countUp()
	for i := 1; i <= 5; i++ {
		v := f4(2)
		fmt.Printf("%v\n", v)
	}

	myf := countUp()
	myg := countUp()
	fmt.Println(myf(2))
	fmt.Println(myf(2))
	fmt.Println(myg(2))
	fmt.Println(myg(2))
}
