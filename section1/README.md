# セクション1:はじめに

## Go言語の特徴
- 静的型付け言語
- コンパイラ言語
- 並行処理の実装が容易
- 実行速度が速い
- GCの遅延が許容できない用途には向かない
  - そのようなケースではC, C++, Rustなど
  
## 保存時の自動整形
- vscodeでsettings.jsonに以下を追記することでファイル保存時に自動整形できる
```
 "[go]": {
        "editor.defaultFormatter": "golang.go",
        "editor.formatOnSave": true,
    },
```