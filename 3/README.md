# セクション3: Go言語・並行処理


<!-- vscode-markdown-toc -->
* 1. [構成](#)
* 2. [メモ](#-1)
	* 2.1. [ロジカルコアとフィジカルコア](#-1)
	* 2.2. [並列処理と並行処理](#-1)
	* 2.3. [runtime schedulerの仕組み](#runtimescheduler)
		* 2.3.1. [Preemption](#Preemption)
		* 2.3.2. [Work stealing](#Workstealing)
		* 2.3.3. [Handoff](#Handoff)
	* 2.4. [Fork join model](#Forkjoinmodel)
	* 2.5. [chained methodのdefer](#chainedmethoddefer)
	* 2.6. [trace](#trace)
	* 2.7. [goroutineの注意点](#goroutine)
	* 2.8. [channel](#channel)
		* 2.8.1. [バッファなし](#-1)
		* 2.8.2. [バッファあり](#-1)
	* 2.9. [goroutine leak](#goroutineleak)
	* 2.10. [channelのclose](#channelclose)
		* 2.10.1. [バッファなしchannelのclose](#channelclose-1)
		* 2.10.2. [バッファなしchannelのclose](#channelclose-1)
	* 2.11. [channelのカプセル化](#channel-1)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

##  1. <a name=''></a>構成
- 00-goroutine : tracerとsyncGroupを使った並列処理入門
- 01-channel-0 : バッファありなしchannel, goroutine leakについて
- 02-channel-1 : channelのclose, カプセル化, 通知専用channelについて
- 03-select-0  : select, timeout contextについて 

##  2. <a name='-1'></a>メモ
###  2.1. <a name='-1'></a>ロジカルコアとフィジカルコア
- ロジカルコア : 論理コア
  - コンピュータから見たときに存在するように見えるCPUのコア
- フィジカルコア : 物理コア
  - CPUの中に実際にあるコア
- 8コア/16スレッドなどのCPUの場合は物理コア数が8で論理コア数が16ということ

###  2.2. <a name='-1'></a>並列処理と並行処理
- 並列処理(Parallelism)
  - 複数コアを使って複数のタスクをそれぞれ異なるコアに格納して同時に実行する方式
- 並行処理(Concurrency)
  - 1つのCPUコアが複数のプロセスを切り替えながら実行することで実際には同時実行していないが、同時実行しているようにあたかも見えるというもの
###  2.3. <a name='runtimescheduler'></a>runtime schedulerの仕組み
- P個のlogical processorを持つ
- 各logical processorはlocal queueを1つもつ
  - このqueueに最大256個のgoroutineを格納できる
- local queueがmaxに達するとglobal queueに入れる
  - global queueはすべてのlogical processorで共有されるもの
- local processorはlocal queueからgoroutineを持ってきてOSのthreadに割り当てて実行させる
- これによりロジカルコア数以上のgoroutineを扱うことができる
- 以下の3つの処理によって割り当てが最適化されている
####  2.3.1. <a name='Preemption'></a>Preemption
- OSのthreadで10ms実行したgoroutineはglobal queueの末尾に移動する
- これによって1つのgoroutineが1つのthreadを長時間占有するということを防ぐことができる
####  2.3.2. <a name='Workstealing'></a>Work stealing
- logical processorは一定間隔ごとに自身のlocal queueを確認してgoroutineがあるならthreadに割り当てる
- また一定間隔ごとにglobal queueを確認してgoroutineがあるなら持ってくる
- そして、自身のlocal queueとglobal queueの両方が空の場合に他のlocal queueに格納されているgoroutineの半分を自身のlocal queueに持ってくる(これがwork stealing)
####  2.3.3. <a name='Handoff'></a>Handoff
- あるスレッドで実行されているgoroutineが待ち状態になった時にlogical processorとそのthreadを切り離して、logical processorに別のthreadを割り当てるというもの
  - これによって待ち状態になっているgoroutineの実行を待つことなく自身のqueueにあるgoroutineをもう一つthreadに割り当てて実行できる

###  2.4. <a name='Forkjoinmodel'></a>Fork join model
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

###  2.5. <a name='chainedmethoddefer'></a>chained methodのdefer
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

###  2.6. <a name='trace'></a>trace
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

###  2.7. <a name='goroutine'></a>goroutineの注意点
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

###  2.8. <a name='channel'></a>channel
- 2つの別のgoroutineで値をやり取りするためのパイプのようなもの
- チャネル変数chに対して`ch<-`で書き込み、`<-ch`で読み込みとなる
- バッファなしとバッファありがある
####  2.8.1. <a name='-1'></a>バッファなし
- 送信側(ch<-)は受信側(<-ch)がくるまで実行停止する、逆もしかり
- プログラマが明示的に何もしなくても送信側受信側で同期がとれるまで停止してくれるから便利
```go
ch := make(chan int)
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	ch <- 10
	time.Sleep(500 * time.Millisecond)
}()
fmt.Println(<-ch)
wg.Wait()
```
####  2.8.2. <a name='-1'></a>バッファあり
- バッファなしの場合はチャネルの受信が開始している状態じゃないと書き込みができなかったが、バッファつきチャネルでは容量があるならok
```go
ch2 := make(chan int, 1)
ch2 <- 2
fmt.Println(<-ch2)
```
- makeの第二引数でバッファサイズを指定
- このコードのように受信が開始していない状態でもバッファがあれば書き込みができる
- 以下のようにバッファサイズを超えるとデッドロックになる
```go
ch2 := make(chan int, 1)
ch2 <- 2
ch2 <- 2
fmt.Println(<-ch2)
```
- またバッファつきチャネルはqueue
```go
ch := make(chan int, 2)
ch <- 1
ch <- 2 
fmt.Println(<-ch)
fmt.Println(<-ch)
// 出力結果 : 1 2
```

###  2.9. <a name='goroutineleak'></a>goroutine leak
- 以下のように無限に待ち状態になって停止していてメモリの解放されないgoroutineがあることをgoroutine leakという
```go
func main(){
    ch1 := make(chan int)
    go func() { // goroutine leak
    	fmt.Println(<-ch1)
    }()
}
```
- goleakを検出するテストは"go.uber.org/goleak"を使って以下のようにかける
```go
func TestLeak(t *testing.T) {
	defer goleak.VerifyNone(t) // ここまで固定
    // この下にgoleakがあるかどうかを検出したいgoroutine対象を複数記述
	main()
}
```
- これでgoleakがあるとテスト結果がFAILになる

###  2.10. <a name='channelclose'></a>channelのclose
####  2.10.1. <a name='channelclose-1'></a>バッファなしchannelのclose
```go
ch1 := make(chan int)
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	fmt.Println(<-ch1)
}()
ch1 <- 10
close(ch1)
v, ok := <-ch1
fmt.Printf("%v %v\n", v, ok) // 0 false
```
- このようにcloseした後に読み込み操作をするとcloseしていなければデッドロックとなってしまうが(0, false)が返却される
  - 0はint型のデフォルト値
####  2.10.2. <a name='channelclose-1'></a>バッファなしchannelのclose
```go
ch2 := make(chan int, 2)
ch2 <- 1
ch2 <- 2
close(ch2)
v, ok = <-ch2
fmt.Printf("%v %v\n", v, ok) // 1 true
v, ok = <-ch2
fmt.Printf("%v %v\n", v, ok) // 2 true
v, ok = <-ch2
fmt.Printf("%v %v\n", v, ok) // 0 false
```
- このようにcloseしたあともバッファに値が残っている間は値が返却される
- バッファに値が残っていない状態になると型のデフォルト値が返って返り値の2つ目はfalseになる
- これを利用して以下のように通知専用のチャネルを作成できる
　- <-chによって待ち状態になっているものがcloseによって解放されるため
```go
// 通知専用のchannelの例
nCh := make(chan struct{}) // 空の構造体のサイズは0byte
for i := 0; i < 5; i++ {
	wg.Add(1)
	go func(i int) {
		defer wg.Done()
		fmt.Printf("goroutine %v started\n", i)
		<-nCh // ここで停止して待ち状態に入る
		fmt.Println(i)
	}(i)
}
time.Sleep(2 * time.Second)
close(nCh) // これにより解放される
fmt.Println("unlocked by manual close")
wg.Wait()
fmt.Println("finish")
```
- 実行結果は以下
```
goroutine 0 started
goroutine 4 started
goroutine 3 started
goroutine 1 started
goroutine 2 started
unlocked by manual close
1
2
3
0
4
```
###  2.11. <a name='channel-1'></a>channelのカプセル化
- 以下のようにカプセル化できる
```go
func main(){
	ch3 := generateCountStream()
	for v := range ch3 {
		fmt.Println(v)
	}
}
// <-chanは読み取り専用のchannel型
func generateCountStream() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i <= 5; i++ {
			ch <- i
		}
	}()
	return ch
}
// 実行結果 0 1 2 3 4 5
```

### select
- channelに値が入っていない場合に受信はブロックするが、ブロックせずに処理を行いたい場合に使える
```go
select {
	case v := <- ch:
		fmt.Println(v)
	default:
		fmt.Println("no value")
}
```
- chに値が入っている場合は`case v := <- ch`が実行されて値が入っていない場合はdefaultが実行される
- このように値が入っていない場合に実行停止せずにほかの処理を走らせたいときにselectを使える
- また、channelがcloseされているとcase文が実行される(default値, false)が返ってくるから
- caseの判定は上から順に判定されて最初にtrueのものが実行される
  - switch文のようなもの
  - つまり複数のcaseが条件を満たしていても1つのcaseしか実行されない
```go
ch1 := make(chan string)
ch2 := make(chan string)
var wg sync.WaitGroup
wg.Add(2)
go func() {
	defer wg.Done()
	time.Sleep(500 * time.Millisecond)
	ch1 <- "A"
}()
go func() {
	defer wg.Done()
	time.Sleep(800 * time.Millisecond)
	ch2 <- "B"
}()
cnt := 0
for ch1 != nil || ch2 != nil {
	cnt++
	select {
	case v := <-ch1:
		fmt.Println(v)
		ch1 = nil
	case v := <-ch2:
		fmt.Println(v)
		ch2 = nil
	}
}
wg.Wait()
fmt.Println("finish")
fmt.Println(cnt)
// 実行結果
// A
// B
// finish
// 2
```
### context.WithTimeout
- タイムアウトの設定方法
```go
ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
defer cancel()
select{
	case <-ctx.Done():
		fmt.Println("Timeout")
	case ...
	...
}
```