# セクション3: Go言語・並行処理

## 構成

## メモ
### ロジカルコアとフィジカルコア
- ロジカルコア : 論理コア
  - コンピュータから見たときに存在するように見えるCPUのコア
- フィジカルコア : 物理コア
  - CPUの中に実際にあるコア
- 8コア/16スレッドなどのCPUの場合は物理コア数が8で論理コア数が16ということ

### 並列処理と並行処理
- 並列処理(Parallelism)
  - 複数コアを使って複数のタスクをそれぞれ異なるコアに格納して同時に実行する方式
- 並行処理(Concurrency)
  - 1つのCPUコアが複数のプロセスを切り替えながら実行することで実際には同時実行していないが、同時実行しているようにあたかも見えるというもの
### runtime schedulerの仕組み
- P個のlogical processorを持つ
- 各logical processorはlocal queueを1つもつ
  - このqueueに最大256個のgoroutineを格納できる
- local queueがmaxに達するとglobal queueに入れる
  - global queueはすべてのlogical processorで共有されるもの
- local processorはlocal queueからgoroutineを持ってきてOSのthreadに割り当てて実行させる
- これによりロジカルコア数以上のgoroutineを扱うことができる
- 以下の3つの処理によって割り当てが最適化されている
#### Preemption
- OSのthreadで10ms実行したgoroutineはglobal queueの末尾に移動する
- これによって1つのgoroutineが1つのthreadを長時間占有するということを防ぐことができる
#### Work stealing
- logical processorは一定間隔ごとに自身のlocal queueを確認してgoroutineがあるならthreadに割り当てる
- また一定間隔ごとにglobal queueを確認してgoroutineがあるなら持ってくる
- そして、自身のlocal queueとglobal queueの両方が空の場合に他のlocal queueに格納されているgoroutineの半分を自身のlocal queueに持ってくる(これがwork stealing)
#### Handoff
- あるスレッドで実行されているgoroutineが待ち状態になった時にlogical processorとそのthreadを切り離して、logical processorに別のthreadを割り当てるというもの
  - これによって待ち状態になっているgoroutineの実行を待つことなく自身のqueueにあるgoroutineをもう一つthreadに割り当てて実行できる

### Fork join model
- main関数もgoroutine
- main関数で新しいgoroutineを走らせるとmainの実行が終わるとほかのgoroutineの終了を待たずに終了してしまう
- 他のgoroutineがすべて終了してからmainのgoroutineを終了させないといけない
- これがjoin
- コードだと以下
```go
var wg sync.WaitGroup // sync.WaitGroupで管理
	wg.Add(1) // カウンタが1増える
	go func() { // goroutine
		defer wg.Done() // Done()でカウンタが1減る(deferを付けているのでこのgoroutineが終了するとカウンタ1減る)
		fmt.Println("goroutine invoked")
	}()
	wg.Wait() // カウンタが0になるまでここで停止する
```

### chained methodのdefer
```go
defer trace.StartRegion(ctx, name).End() // chained methodなのでEnd()だけdefer
```
- これは以下と同じ
```go
p := trace.StartRegion(ctx, name)
defer p.End()
```
- chained methodをdeferすると最後のメソッドだけdeferされてそれ以外は即時実行される
- 極端な例は以下
  - https://stackoverflow.com/questions/68437177/what-will-happen-when-using-defer-in-chained-method-call
```go
defer A().B().C().D().E().F().G().H()
// Hだけdeferされる、それ以外はこのdefer文に到達したときに普通に実行される
```

### trace
- goroutineがどう動いているかを可視化することができる
- 以下使用例(逐次処理の場合)
```go
// task定義例
func task(ctx context.Context, name string) {
	defer trace.StartRegion(ctx, name).End() // chained methodなのでEnd()だけdefer
	time.Sleep(time.Second)
	fmt.Println(name)
}

```
```go
// 使用例
// ログ用のファイル作成
f, err := os.Create("trace.out")
if err != nil {
	log.Fatalln("Error:", err)
}
defer func() {
	if err := f.Close(); err != nil {
		log.Fatalln("Error:", err)
	}
}()
// trace利用 trace.Start()の引数にファイルを与える
if err := trace.Start(f); err != nil {
	log.Fatalln("Error:", err)
}
defer trace.Stop() // deferを使って止める
ctx, t := trace.NewTask(context.Background(), "main") // 第二引数は任意のタスク名
defer t.End()
fmt.Println("The number of logical CPU Cores:", runtime.NumCPU())
task(ctx, "Task1")
task(ctx, "Task2")
task(ctx, "Task3")
```
- 実行終了後に`go tool trace ファイル名`を実行するとブラウザが立ち上がる
  - この例では`go tool trace trace.out`
- ブラウザのUser-defined tasksをクリック -> Countの数字をクリック -> EventsのTask1をクリック
- そうするとそれぞれのtaskがどのように実行されたか可視化されたものが表示される
- この例では逐次実行なので、task1->task2->task3->が順に実行されている

### goroutineの注意点
```go
s := []int{1, 2, 3}
for _, i := range s {
    wg.Add(1)
	go func() {
    	defer wg.Done()
		fmt.Println(i)
	}()
}
// 実行結果 : 3 3 3
```
- goroutineの実行開始に時間がかかり、goroutineが実行開始してiにアクセスしたころには次のfor文にいっているためにすべて3になりうる
- これを回避するためにはgoroutineの内部から外部の変数にアクセスしないように引数で渡せばよい
```go
s := []int{1, 2, 3}
for _, i := range s {
    wg.Add(1)
	go func(i int) {
    	defer wg.Done()
		fmt.Println(i)
	}(i)
}
// 実行結果 : 1 3 2
```
- これで先ほどの問題を解決できるがgoroutineがそれぞれ並列処理で実行されるため順番に実行されるとは限らない