# Udemy(Go言語の基礎と並行処理)

<!-- vscode-markdown-toc -->
* 1. [ソース](#)
* 2. [構成](#-1)
* 3. [自分用メモまとめ](#-1)
	* 3.1. [Go言語の特徴](#Go)
	* 3.2. [保存時の自動整形](#-1)
	* 3.3. [module・packageまわり](#modulepackage)
	* 3.4. [外部モジュール利用法](#-1)
	* 3.5. [変数宣言](#-1)
	* 3.6. [シャドーイング](#-1)
	* 3.7. [配列とスライス](#-1)
	* 3.8. [メソッドとレシーバ](#-1)
	* 3.9. [defer](#defer)
	* 3.10. [無名関数](#-1)
	* 3.11. [closure](#closure)
	* 3.12. [interface](#interface)
	* 3.13. [range](#range)
	* 3.14. [error](#error)
	* 3.15. [generics](#generics)
	* 3.16. [ユニットテスト](#-1)
	* 3.17. [logger](#logger)
	* 3.18. [ロジカルコアとフィジカルコア](#-1)
	* 3.19. [並列処理と並行処理](#-1)
	* 3.20. [runtime schedulerの仕組み](#runtimescheduler)
		* 3.20.1. [Preemption](#Preemption)
		* 3.20.2. [Work stealing](#Workstealing)
		* 3.20.3. [Handoff](#Handoff)
	* 3.21. [Fork join model](#Forkjoinmodel)
	* 3.22. [chained methodのdefer](#chainedmethoddefer)
	* 3.23. [trace](#trace)
	* 3.24. [goroutineの注意点](#goroutine)
	* 3.25. [channel](#channel)
		* 3.25.1. [バッファなし](#-1)
		* 3.25.2. [バッファあり](#-1)
	* 3.26. [goroutine leak](#goroutineleak)
	* 3.27. [channelのclose](#channelclose)
		* 3.27.1. [バッファなしchannelのclose](#channelclose-1)
		* 3.27.2. [バッファなしchannelのclose](#channelclose-1)
	* 3.28. [channelのカプセル化](#channel-1)
	* 3.29. [select](#select)
	* 3.30. [context.WithTimeout](#context.WithTimeout)
	* 3.31. [data race](#datarace)
	* 3.32. [Mutex](#Mutex)
	* 3.33. [RWMutex](#RWMutex)
	* 3.34. [atomic](#atomic)
	* 3.35. [Context](#Context)
	* 3.36. [errgroup](#errgroup)
	* 3.37. [pipeline](#pipeline)
	* 3.38. [fan-out, fan-in](#fan-outfan-in)
	* 3.39. [heartbeat, watchdog Timer](#heartbeatwatchdogTimer)
	* 3.40. [selectのランダム性について](#select-1)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

##  1. <a name=''></a>ソース
- Go言語の基礎と並行処理
- https://www.udemy.com/course/golang_concurrency/?couponCode=KEEPLEARNING

##  2. <a name='-1'></a>構成
- 1 : セクション1(はじめに)
- 2 : セクション2(Go言語・基礎)
- 3 : セクション3(Go言語・並行処理)

##  3. <a name='-1'></a>自分用メモまとめ
###  3.1. <a name='Go'></a>Go言語の特徴
- 静的型付け言語
- コンパイラ言語
- 並行処理の実装が容易
- 実行速度が速い
- GCの遅延が許容できない用途には向かない
  - そのようなケースではC, C++, Rustなど
  
###  3.2. <a name='-1'></a>保存時の自動整形
- vscodeでsettings.jsonに以下を追記することでファイル保存時に自動整形できる
```
 "[go]": {
        "editor.defaultFormatter": "golang.go",
        "editor.formatOnSave": true,
    },
```
###  3.3. <a name='modulepackage'></a>module・packageまわり
- package : 同じディレクトリに存在するソースコードファイル群のこと
- module : パッケージの集合
- go.mod : Goモジュールのパスを書いておくファイル
  - 一番最初に`go mod init "名前"`で作成する
- go.sum : 依存モジュールのチェックサムが記録されているファイル
  - チェックサムとはファイルやデータの一意性を示すハッシュのこと
  - ファイルが改善されていないことを確認できる
  - このチェックサムを用いることでモジュールの内容に変更があったかどうかを検出
- 変数や関数を外部パッケージからアクセス可能にするためには先頭を大文字にしないといけない
- 先頭小文字のものは同じパッケージ内の別ファイルからアクセス可能
- Javaの4つのアクセスレベルのうちpublic, package private(アクセス修飾子ないやつ)にGoの大文字小文字が対応
###  3.4. <a name='-1'></a>外部モジュール利用法
- 以下の3段階(例としてgodotenvを利用する場合)
    1. `go get github.com/joho/godotenv`
    2. import記述(`import "github.com/joho/godotenv"`)
    3. 利用(`godotenv.メソッド名()`など)
- VSCodeの拡張機能を使うと以下のように利用可能
    1. 利用(`godotenv.メソッド名()`)
    2. 自動でimportが追加されるが波線が引かれる(go getしていないから)
    3. `go mod tidy`
        - これはimportしていてまだgo getしていないモジュールをすべてgetしてgo.mod, go.sumに記述するコマンド

###  3.5. <a name='-1'></a>変数宣言
- `var i int`
  - 明示的な代入をスキップ可能(自動初期化)
  - 自動初期化は以下
    - bool : false
    - string : ""
    - pointer : nil
    - int : 0
- `i := 1`
  - 明示的な代入が必須
  - ローカル変数として宣言されたブロック内が有効なスコープ
  - そのため関数外では使用不可(グローバル変数になれない)

###  3.6. <a name='-1'></a>シャドーイング
- 同じスコープ内で同じ名前の変数を再宣言することによって、外側の変数を隠す効果を持つ機能

###  3.7. <a name='-1'></a>配列とスライス
- 配列は静的でスライスは動的
- スライスは参照型
- 配列とスライスには要素数と容量というものがあり、len, capで取得できる
  - 要素数 : 要素の個数
  - 容量 : 確保されているメモリの個数
- スライスでappendを使って容量オーバーすると現在の容量の2倍のメモリを確保した新しい領域に値が移動して、現在の変数の参照先アドレスが変わる
- 以下の記事が非常に参考になるし大事
- https://qiita.com/Kashiwara/items/e621a4ad8ec00974f025
- メモリまわりの振る舞いについては03-slice-map/main.goのコメントアウト参照

###  3.8. <a name='-1'></a>メソッドとレシーバ
- メソッドとは、レシーバを持つ関数のこと
- 以下の例だと`task Task`がレシーバ
- `func (レシーバ変数名 その変数の構造体の型) 関数名(引数)返り値{}`
```go
func (task Task) extendEstimate() {
	task.Estimate += 10
}
```
- 使い方は普通の言語のオブジェクトと同じ
```go
t := Task{}
t.extendEstimate()
```
- 値レシーバ : レシーバが構造体型
  - 構造体をコピーして呼び出すため、元のオブジェクトのフィールドの値は変更されない
- ポインタレシーバ : レシーバが構造体のポインタ型
  - 構造体のポインタが渡されるので元のオブジェクトのフィールドの値を変更できる
- 詳しくは04-struct-receiver/main.go参照

###  3.9. <a name='defer'></a>defer
- deferから始まる文はその文が宣言されている関数の実行終了直前で呼び出される
- deferを複数利用するとstackで扱われる
  - 逆順に実行されるということ
- またdeferのついた関数(遅延関数)の引数に与えた変数は即時評価される
  - defer文の後で引数の値を変更しても影響受けない
  - 詳しくは以下
    - https://qiita.com/Ishidall/items/8dd663de5755a15e84f2
###  3.10. <a name='-1'></a>無名関数
- 即時実行する場合は以下のように無名関数宣言後の`}`に続けて`()`を付けて引数を与えて呼び出す
```go
func(i int) {
	fmt.Println(i)
}(i)
```
- 変数に代入して実行する場合は以下のように関数型の変数を利用できる
```go
f1 := func(i int) int {
	return i + 1
}
fmt.Println(f1(i))
```
- 引数で無名関数を利用する例は以下
```go
func addExt(f func(file string) string, name string) {
	fmt.Println(f(name))
}
func main(){
    f2 := func(file string) string {
		return file + ".csv"
	}
	// 引数で無名関数を利用する例
	addExt(f2, "file1")
}
```
- 返り値で無名関数を利用する例は以下
```go
func multiply() func(int) int {
	return func(n int) int {
		return n * 1000
	}
}
func main(){
    f3 := multiply()
	fmt.Println(f3(2))
}
```
###  3.11. <a name='closure'></a>closure
- 関数が実行されたときにその性的スコープで定義された変数を利用できる関数のこと
```go
// closure
func countUp() func(int) int {
	count := 0
	return func(n int) int {
		count += n
		return count
	}
}
func main(){
    f4 := countUp()
	for i := 1; i <= 5; i++ {
		v := f4(2)
		fmt.Printf("%v\n", v) // 2, 4, 6, 8, 10
	}
}
```
- f4という関数変数がcountUp()の関数内で定義された変数と、返り値の無名関数を持っているとイメージするとわかりやすい
- 変数がcountの値を持っているというイメージは以下を考えるとわかりやすい
```go
f := countUp()
g := countUp()
fmt.Println(f(2)) // 2
fmt.Println(f(2)) // 4
fmt.Println(g(2)) // 2
fmt.Println(g(2)) // 4
```
- 詳しくは以下参照
  - https://blog.framinal.life/entry/2022/12/12/090012

###  3.12. <a name='interface'></a>interface
- Javaと同じメソッドをまとめたもの
- interfaceで定義されているメソッドと同名のメソッドを実装しているのであればどんな型でも扱える
- Javaではimplementsでinterfaceを明示する必要があるが、Goは不必要
  - interfaceで定義されているメソッドをすべて実装している構造体はそのinterfaceをimplementsしていると勝手に解釈してくれる
- interfaceの定義例は以下
```go
type controller interface {
	speedUp() int
	speedDown() int
}
```
- 利用例は以下
```go
func speedUpAndDown(c controller) {
	fmt.Printf("current speed: %v\n", c.speedUp())
	fmt.Printf("current speed: %v\n", c.speedDown())
}
```
- また、`interface{}`は`any`と同じ意味を持つ
  - `interface`はそのブロック内で定義されているメソッドをすべて実装している任意の型を扱える
  - `interface{}`は実装する必要のあるメソッドが何もなく制約が何もないと考えると`any`と同値であると理解できる
  - APIのレスポンスで見る`map[string]interface{}`はstringがkeyでvalueはintでもstringでも何でもよいよということ

###  3.13. <a name='range'></a>range
- `for`文で`range`を使うことでスライスの(index, element)を取得できる
- ただ注意点として、elementは値のコピーであって参照をコピーしていないので元の値を変えることはできない
```go
type item struct {
	price float32
}
items := []item{
	{price: 10.},
	{price: 20.},
	{price: 30.},
}
for _, item := range items {
	item.price *= 1.1 // このitemは構造体の値のコピーであって参照ではないので値変わらない
}
fmt.Printf("%+v\n", items) // priceは変わっていない
for i := range items {
	items[i].price *= 1.1 // 値を変えたい場合はこのようにindexを使ってアクセス
}
fmt.Printf("%+v\n", items)
```

###  3.14. <a name='error'></a>error
- 08-errors/main.goのコメントアウトや使用例を参照

###  3.15. <a name='generics'></a>generics
- C++のtemplateみたいなもの
- 以下のように"func 関数名\[関数内で使う型変数 型\]"で"\[\]"の中で定義して使える
```go
// int型またはfloat64型を使う場合 int | float64とかく
func Sum[T int | float64](nums []T) T {
	var sum T
	for _, v := range nums {
		sum += v
	}
	return sum
}
```
- また型の列挙にはinterfaceを使える
```go
type customConstraints interface {
	int | int16 | float32 | float64 | string
}
func add[T customConstraints](x, y T) T {
	return x + y
}
```
- ある型から派生させた独自の型も対象にしたいときはチルダを付ける
```go
type NewInt int // int型から派生した独自のNewInt型
type customConstraints interface {
	~int | int16 | float32 | float64 | string // ~intとすることでint型から派生した独自の型も対象
}
func add[T customConstraints](x, y T) T {
	return x + y
}
var i1, i2 NewInt = 3, 4
fmt.Println(add(i1, i2))
```
- genericsは関数内で複数定義可能であるため、以下のように利用できる
```go
func sumValues[K int | string, V constraints.Float | constraints.Integer](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}
m1 := map[string]uint{
	"A": 1,
	"B": 2,
	"C": 3,
}
m2 := map[int]float32{
	1: 1.23,
	2: 4.56,
	3: 7.89,
}
fmt.Println(sumValues(m1))
fmt.Println(sumValues(m2))
```

###  3.16. <a name='-1'></a>ユニットテスト
- vscodeのgoの拡張機能を入れている場合、関数を右クリックしてGo : generate unit tests for functionを選択すれば自動でひな形が生成される
- 自動で以下が生成される
```go
func TestAdd(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
```
- これに対してテストケースを追加していくだけ
```go
func TestAdd(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "1+2=3",          // テストケース名
			args: args{x: 1, y: 2}, // 引数
			want: 3,                // 期待する返り値
		},
		{
			name: "2+2=4",
			args: args{x: 2, y: 2},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
```
- ユニットテストの実行は`go test -v .`
  - vオプションは詳細表示という意味
- coverageの確認は`go test -v -cover -coverprofile=coverage.out .`
  - coverageとはテスト対象の関数のうちテストできている箇所の割合のこと
  - すべての分岐ルートをたどれているかということ
- coverageが100%でないときに、テストできていないソースコード箇所を見つけるコマンドは`go tool cover -html=coverage.out`
  - ブラウザが立ち上がり、テストできていないソースコード該当箇所が赤くなる
  
###  3.17. <a name='logger'></a>logger
- 以下のようにlog.New()で作れる
- 第一引数には出力先を指定、io.MultiWriterで複数指定が可能
  - ここでは指定したfileと標準エラー出力を指定
- 第3引数のlog.Lshortfileはファイル名や行数などを表示してくれる
```go
file, err := os.Create("log.txt")
warnLogger := log.New(io.MultiWriter(file, os.Stderr), "WARN: ", log.Lshortfile)
warnLogger.Println("warning A") // WARN: main.go:20: warning A
```
- PrintlnではなくFatallnを使うとログを出力した後にプログラムを強制終了させることができる



###  3.18. <a name='-1'></a>ロジカルコアとフィジカルコア
- ロジカルコア : 論理コア
  - コンピュータから見たときに存在するように見えるCPUのコア
- フィジカルコア : 物理コア
  - CPUの中に実際にあるコア
- 8コア/16スレッドなどのCPUの場合は物理コア数が8で論理コア数が16ということ

###  3.19. <a name='-1'></a>並列処理と並行処理
- 並列処理(Parallelism)
  - 複数コアを使って複数のタスクをそれぞれ異なるコアに格納して同時に実行する方式
- 並行処理(Concurrency)
  - 1つのCPUコアが複数のプロセスを切り替えながら実行することで実際には同時実行していないが、同時実行しているようにあたかも見えるというもの
###  3.20. <a name='runtimescheduler'></a>runtime schedulerの仕組み
- P個のlogical processorを持つ
- 各logical processorはlocal queueを1つもつ
  - このqueueに最大256個のgoroutineを格納できる
- local queueがmaxに達するとglobal queueに入れる
  - global queueはすべてのlogical processorで共有されるもの
- local processorはlocal queueからgoroutineを持ってきてOSのthreadに割り当てて実行させる
- これによりロジカルコア数以上のgoroutineを扱うことができる
- 以下の3つの処理によって割り当てが最適化されている
####  3.20.1. <a name='Preemption'></a>Preemption
- OSのthreadで10ms実行したgoroutineはglobal queueの末尾に移動する
- これによって1つのgoroutineが1つのthreadを長時間占有するということを防ぐことができる
####  3.20.2. <a name='Workstealing'></a>Work stealing
- logical processorは一定間隔ごとに自身のlocal queueを確認してgoroutineがあるならthreadに割り当てる
- また一定間隔ごとにglobal queueを確認してgoroutineがあるなら持ってくる
- そして、自身のlocal queueとglobal queueの両方が空の場合に他のlocal queueに格納されているgoroutineの半分を自身のlocal queueに持ってくる(これがwork stealing)
####  3.20.3. <a name='Handoff'></a>Handoff
- あるスレッドで実行されているgoroutineが待ち状態になった時にlogical processorとそのthreadを切り離して、logical processorに別のthreadを割り当てるというもの
  - これによって待ち状態になっているgoroutineの実行を待つことなく自身のqueueにあるgoroutineをもう一つthreadに割り当てて実行できる

###  3.21. <a name='Forkjoinmodel'></a>Fork join model
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

###  3.22. <a name='chainedmethoddefer'></a>chained methodのdefer
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

###  3.23. <a name='trace'></a>trace
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

###  3.24. <a name='goroutine'></a>goroutineの注意点
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

###  3.25. <a name='channel'></a>channel
- 2つの別のgoroutineで値をやり取りするためのパイプのようなもの
- チャネル変数chに対して`ch<-`で書き込み、`<-ch`で読み込みとなる
- バッファなしとバッファありがある
####  3.25.1. <a name='-1'></a>バッファなし
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
####  3.25.2. <a name='-1'></a>バッファあり
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

###  3.26. <a name='goroutineleak'></a>goroutine leak
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

###  3.27. <a name='channelclose'></a>channelのclose
####  3.27.1. <a name='channelclose-1'></a>バッファなしchannelのclose
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
####  3.27.2. <a name='channelclose-1'></a>バッファなしchannelのclose
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
###  3.28. <a name='channel-1'></a>channelのカプセル化
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

###  3.29. <a name='select'></a>select
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
- caseを途中で抜ける場合はbreakを使う
```go
select{
	case ...:
		break
		...
}
```
###  3.30. <a name='context.WithTimeout'></a>context.WithTimeout
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


###  3.31. <a name='datarace'></a>data race
- データレースとは、あるメモリー位置への書き込みであって、その同じ位置に対するほかの読み込みまたは書き込みと並列に起きるもの
- 以下がデータレースの例
```go
var wg sync.WaitGroup
var i int
wg.Add(2)
go func() {
	defer wg.Done()
	i++
}()
go func() {
	defer wg.Done()
	i++
}()
wg.Wait()
fmt.Println(i)
```
- 2つのgoroutineがたまたま同時にiのインクリメント(0から1)を行うと実行結果が2ではなく1になりうる
- これがデータレース
- データレースの検出方法は`go run -race main.go`
- このようなデータ競合の回避にmutexを使える
###  3.32. <a name='Mutex'></a>Mutex
- goの排他制御機構
```go
var wg sync.WaitGroup
var mu sync.Mutex
var i int
wg.Add(2)
go func() {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	i++
}()
go func() {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	i++
}()
wg.Wait()
fmt.Println(i)
```
- Lock()でMutexをロックしてUnlock()でMutexをアンロックする
- これによってiに対するdata raceを回避できている
  - iを占有できているから

###  3.33. <a name='RWMutex'></a>RWMutex
- readのlockによって他のreadをlockしない
```go
func main(){
	var wg sync.WaitGroup
	var rwMu sync.RWMutex
	var c int
	wg.Add(4)
	go write(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)
	wg.Wait()
	fmt.Println("finish")
}
func read(mu *sync.RWMutex, wg *sync.WaitGroup, c *int) {
	defer wg.Done()
	time.Sleep(10 * time.Millisecond)
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("read lock")
	fmt.Println(*c)
	time.Sleep(1 * time.Second)
	fmt.Println("read unlock")
}
func write(mu *sync.RWMutex, wg *sync.WaitGroup, c *int) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("write lock")
	*c += 1
	time.Sleep(1 * time.Second)
	fmt.Println("write unlock")
}
// 実行結果
// write lock
// write unlock
// read lock
// 1
// read lock
// read lock
// 1
// 1
// read unlock
// read unlock
// read unlock
// finish
```
###  3.34. <a name='atomic'></a>atomic
- 簡単に排他制御できる
- 以下の例だと第一引数にロックする対象の整数変数のポインタ、第二引数に加算する値
```go
var wg sync.WaitGroup
var c int64
for i := 0; i < 5; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 10; j++ {
			atomic.AddInt64(&c, 1) // 排他制御
		}
	}()
}
wg.Wait()
fmt.Println(c)
fmt.Println("finish")
// 実行結果
// 50
```

###  3.35. <a name='Context'></a>Context
- main goroutineからサブgoroutinに情報を伝搬させるときにつかう
- goroutineの関係が木構造になっていて親から子に情報を伝搬させていくということ
- 主な機能例は以下
  - `func WichCancel(parent Context)` : マニュアルキャンセル
    - 詳しくは`07-context/withcancel/main.go`参照
  - `func WithDeadline(parent context, d time.Time)` : 時刻でdeadlineを指定
    - 詳しくは`07-context/withdeadline/main.go`参照
  - `func WithTimeout(parent Context, timeout time.Duration)` : timeoutを指定
    - 詳しくは`07-context/withtimeout/main.go`参照
- contexを作成するには第一引数に親のコンテキストを与える
  - つまり、機能を追加したい場合は親のコンテキストをもとに新しいコンテキストを生成するということ
  - `ctx, cancel := context.WithCancel(context.Background())`など
  - `context.Background()`は空のcontextのことであり、親のいないcontext(ルートノード)に使う

###  3.36. <a name='errgroup'></a>errgroup
- errgroupでは複数のgoroutineを実行してそれらのうちにエラーがあったときにエラーを知るということを可能にしてくれる
- `go func(){}`ではerrorを検出できない
- errgroupでは`eg.Go(func() error{`のようにerrorを返せるようになっている
```go
func main() {
	eg := new(errgroup.Group)
	s := []string{"task1", "fake1", "task2", "fake2"}
	for _, v := range s {
		task := v
		eg.Go(func() error {
			return doTask(task)
		})
	}
	// errorがあった場合に最初に発生したものだけ返る
	if err := eg.Wait(); err != nil {
		fmt.Printf("error :%v\n", err)
	}
	fmt.Println("finish")
}
func doTask(task string) error {
	if task == "fake1" || task == "fake2" {
		return fmt.Errorf("%v failed", task)
	}
	fmt.Printf("task %v completed\n", task)
	return nil
}

```
- またcontextを使ってどれか1つのgoroutineがerrorを返したときにほかのgoroutineをすべてcancelするなどの処理が可能
  - `eg := new(errgroup.Group)`ではなく`eg, ctx := errgroup.WithContext(context.Background())`などとする
  - errgroup.WithContextは最初にerrorを検出した段階でDone()チャネルをcloseするという仕様になっている
- 以下例
```go
func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	s := []string{"task1", "fake1", "task2", "fake2", "task3"}
	for _, v := range s {
		task := v
		eg.Go(func() error {
			return doTask(ctx, task)
		})
	}
	// errorがあった場合に最初に発生したものだけ返る
	if err := eg.Wait(); err != nil {
		fmt.Printf("error :%v\n", err)
	}
	fmt.Println("finish")
}
func doTask(ctx context.Context, task string) error {
	var t *time.Ticker
	switch task {
	case "fake1":
		t = time.NewTicker(500 * time.Millisecond)
	case "fake2":
		t = time.NewTicker(700 * time.Millisecond)
	default:
		t = time.NewTicker(200 * time.Millisecond)
	}
	select {
	case <-ctx.Done(): // 他のgoroutineでerrorを返しているとcloseになっている
		fmt.Printf("%v cancelled : %v\n", task, ctx.Err())
		return ctx.Err()
	case <-t.C:
		t.Stop()
		if task == "fake1" || task == "fake2" {
			return fmt.Errorf("%v process failed", task)
		}
		fmt.Printf("task %v completed\n", task)
	}
	return nil
}
```

###  3.37. <a name='pipeline'></a>pipeline
- channelを引数と返り値にとることによってpipelineを実現できる
- それぞれの階層で1つずつ入力が流れていくイメージ
- 詳しくは09-pipeline/main.go参照

###  3.38. <a name='fan-outfan-in'></a>fan-out, fan-in
- pipelineにおいてあるステージの計算量が多くボトルネックになっている場合に使う
- そのステージに順に1つずつ送られてくるchannelをさらにgoroutineとして並列処理させる
  - これがfan-out
- 実行後に出力結果である複数のchannelを合流させて再度順に次のステージに入力する
  - これがfan-in

###  3.39. <a name='heartbeatwatchdogTimer'></a>heartbeat, watchdog Timer
- main goroutineがsub goroutineが正常に動いていることを監視するために使う
- heartbeatとは、sub goroutineからmain goroutineに一定間隔で正常に稼働していることを知らせるもの
  - パルス波のようなイメージ
- watchdog timerは設定されたタイムアウト時間に達するとsub goroutineが稼働していないと異常検知するもの
  - watchdog timerは横軸が時刻のグラフを考えると線形に値が増えていき、縦軸の値がタイムアウトに達すると異常を検知
  - ただし、途中でheartbeatが送信されたタイミングで値を0に戻す
- 詳しくは11-heartbeatのコード参照

###  3.40. <a name='select-1'></a>selectのランダム性について
- 複数のcaseが条件を満たしているときに、select文では一様分布に基づいてランダムに選ばれる
```go
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

```
- 実行結果は以下
```
ch2
ch1
ch1
ch2
ch1
ch2
ch1
ch1
ch2
ch1
```