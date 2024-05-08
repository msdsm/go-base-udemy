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