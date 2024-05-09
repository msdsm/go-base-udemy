package main

import (
	"fmt"
	"unsafe"
)

// 構造体の宣言
type Task struct {
	Title    string
	Estimate int
}

func main() {
	// 構造体の利用
	task1 := Task{
		Title:    "Learn Golang",
		Estimate: 3,
	}
	task1.Title = "Learning Go"
	// 型は"パッケージ名"."構造体名"になる
	fmt.Printf("%[1]T %+[1]v %v\n", task1, task1.Title) // main.Task {Title:Learning Go Estimate:3} Learning Go

	// 構造体のコピーは参照渡しではなく値渡しになる
	task2 := task1
	task2.Title = "new"
	fmt.Printf("task1: %v task2: %v\n", task1.Title, task2.Title)

	// 構造体のポインタ
	task1p := &Task{
		Title:    "Learn concurrency",
		Estimate: 2,
	}
	fmt.Printf("task1p: %T %+v %v\n", task1p, *task1p, unsafe.Sizeof(task1p))
	// (*task1p).Titleは以下のように省略できる
	task1p.Title = "Changed"
	fmt.Printf("task1p: %+v\n", *task1p)
	var task2p *Task = task1p
	task2p.Title = "Changed by Task2"
	fmt.Printf("task1: %+v\n", *task1p)
	fmt.Printf("task2: %+v\n", *task2p)

	// 値レシーバの場合はオブジェクトのフィールド変わらない
	task1.extendEstimate()
	fmt.Printf("task1 value receiver: %+v\n", task1.Estimate) // 3
	// ポインタレシーバの場合は変わる
	// (&task1).extendEstimatePointer()は以下のようにdeferenceの省略が可能
	task1.extendEstimatePointer()
	fmt.Printf("task1 value receiver: %+v\n", task1.Estimate) // 13
}

// 値レシーバ
func (task Task) extendEstimate() {
	task.Estimate += 10
}

// ポインタレシーバ
func (taskp *Task) extendEstimatePointer() {
	taskp.Estimate += 10
}
