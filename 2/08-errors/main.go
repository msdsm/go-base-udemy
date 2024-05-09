package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrCustom = errors.New("not found")

func main() {
	// errors.Newで作成可能
	err01 := errors.New("something wrong")
	err02 := errors.New("something wrong")

	// メモリアドレス、型、値の確認
	// errorはerrors.ErrorStringという構造体のポインタ
	fmt.Printf("%[1]p %[1]T %[1]v\n", err01) // 0x188c038 *errors.errorString something wrong
	// Error()がstringを返すが省略可能
	fmt.Println(err01) // something wrong
	// ポインタなのでfalseになる
	fmt.Println(err01 == err02) // false

	// errorのwrap
	// fmt.Errorf()で%wを使うとfmt.wrapErrorのポインタ型になる
	err0 := fmt.Errorf("add info: %w", errors.New("original error"))
	fmt.Printf("%[1]p %[1]T %[1]v\n", err0) // 0x18960a0 *fmt.wrapError add info: original error
	// Unwrapを使ってwrapしているerrorを取り出せる
	fmt.Println(errors.Unwrap(err0))        // original error
	fmt.Printf("%T\n", errors.Unwrap(err0)) // *errors.errorString

	// %vを使ってwrapするとerrors.errorStringのポインタ型、つまり普通のerrorと同じ型になる
	err1 := fmt.Errorf("add info: %v", errors.New("original error"))
	fmt.Println(err1)        // add info: original error
	fmt.Printf("%T\n", err1) // *errors.errorString
	// %wを使ってwrapしていないためunwrapしても何のerrorもない->nilが返る
	fmt.Println(errors.Unwrap(err1)) // <nil>

	// wrapのwrapの例
	// Errから始まるものは標準のerror
	err2 := fmt.Errorf("in repository layer: %w", ErrCustom)
	fmt.Println(err2) // in repository layer: not found
	err2 = fmt.Errorf("in service layer: %w", err2)
	fmt.Println(err2) // in service layer: in repository layer: not found

	// errors.Is(a,b) : aがbをどこかの階層でwrapしているかどうか(nilが返却されるまでunwrapしてくれる)
	if errors.Is(err2, ErrCustom) {
		fmt.Println("matched") // matched
	}

	file := "dummy.txt"
	err3 := fileChecker(file)
	if err3 != nil {
		if errors.Is(err3, os.ErrNotExist) {
			fmt.Printf("%v file not found\n", file) // dummy.txt file not found
		} else {
			fmt.Println("unknown error")
		}
	}

}
func fileChecker(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("in checker: %w", err)
	}
	defer f.Close()
	return nil
}
