# セクション2 : Go言語・基礎


<!-- vscode-markdown-toc -->
* 1. [構成](#)
* 2. [メモ](#-1)
	* 2.1. [module・packageまわり](#modulepackage)
	* 2.2. [外部モジュール利用法](#-1)
	* 2.3. [変数宣言](#-1)

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

### シャドーイング
- 同じスコープ内で同じ名前の変数を再宣言することによって、外側の変数を隠す効果を持つ機能

### 配列とスライス
- 配列は静的でスライスは動的
- スライスは参照型
- 配列とスライスには要素数と容量というものがあり、len, capで取得できる
  - 要素数 : 要素の個数
  - 容量 : 確保されているメモリの個数
- スライスでappendを使って容量オーバーすると現在の容量の2倍のメモリを確保した新しい領域に値が移動して、現在の変数の参照先アドレスが変わる
- 以下の記事が非常に参考になるし大事
- https://qiita.com/Kashiwara/items/e621a4ad8ec00974f025
- メモリまわりの振る舞いについては03-slice-map/main.goのコメントアウト参照

### メソッドとレシーバ
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