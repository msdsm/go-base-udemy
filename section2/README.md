# セクション2 : Go言語・基礎


<!-- vscode-markdown-toc -->
* 1. [構成](#)
* 2. [メモ](#-1)
	* 2.1. [module・packageまわり](#modulepackage)
	* 2.2. [外部モジュール利用法](#-1)
	* 2.3. [変数宣言](#-1)
	* 2.4. [シャドーイング](#-1)
	* 2.5. [配列とスライス](#-1)
	* 2.6. [メソッドとレシーバ](#-1)
	* 2.7. [defer](#defer)
	* 2.8. [無名関数](#-1)
	* 2.9. [closure](#closure)
	* 2.10. [interface](#interface)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->


##  1. <a name=''></a>構成
- 00-module-package : module, packageについて
- 01-variables : 変数について
##  2. <a name='-1'></a>メモ
###  2.1. <a name='modulepackage'></a>module・packageまわり
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
###  2.2. <a name='-1'></a>外部モジュール利用法
- 以下の3段階(例としてgodotenvを利用する場合)
    1. `go get github.com/joho/godotenv`
    2. import記述(`import "github.com/joho/godotenv"`)
    3. 利用(`godotenv.メソッド名()`など)
- VSCodeの拡張機能を使うと以下のように利用可能
    1. 利用(`godotenv.メソッド名()`)
    2. 自動でimportが追加されるが波線が引かれる(go getしていないから)
    3. `go mod tidy`
        - これはimportしていてまだgo getしていないモジュールをすべてgetしてgo.mod, go.sumに記述するコマンド

###  2.3. <a name='-1'></a>変数宣言
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

###  2.4. <a name='-1'></a>シャドーイング
- 同じスコープ内で同じ名前の変数を再宣言することによって、外側の変数を隠す効果を持つ機能

###  2.5. <a name='-1'></a>配列とスライス
- 配列は静的でスライスは動的
- スライスは参照型
- 配列とスライスには要素数と容量というものがあり、len, capで取得できる
  - 要素数 : 要素の個数
  - 容量 : 確保されているメモリの個数
- スライスでappendを使って容量オーバーすると現在の容量の2倍のメモリを確保した新しい領域に値が移動して、現在の変数の参照先アドレスが変わる
- 以下の記事が非常に参考になるし大事
- https://qiita.com/Kashiwara/items/e621a4ad8ec00974f025
- メモリまわりの振る舞いについては03-slice-map/main.goのコメントアウト参照

###  2.6. <a name='-1'></a>メソッドとレシーバ
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

###  2.7. <a name='defer'></a>defer
- deferから始まる文はその文が宣言されている関数の実行終了直前で呼び出される
- deferを複数利用するとstackで扱われる
  - 逆順に実行されるということ
- またdeferのついた関数(遅延関数)の引数に与えた変数は即時評価される
  - defer文の後で引数の値を変更しても影響受けない
  - 詳しくは以下
    - https://qiita.com/Ishidall/items/8dd663de5755a15e84f2
###  2.8. <a name='-1'></a>無名関数
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
###  2.9. <a name='closure'></a>closure
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

###  2.10. <a name='interface'></a>interface
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

### range
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