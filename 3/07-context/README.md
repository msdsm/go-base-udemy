# セクション3-27
- context: cancel + timeout + deadline

## withtimeout
- WithTimeoutの例
- main goroutineから3つのsub goroutineを一斉に終わらせる例

## withcancel
- WithCancelの例
- normalTaskとcriticalTaskというsub goroutineをmain goroutineから走らせる
- criticalTaskにだけtimeoutを設定
- criticalTaskでtimeoutが発生するとその情報をmain goroutineに伝える
- main goroutineはcriticalTaskからそのような情報を受け取ることによってほかのgoroutineをすべて終了させる

### withdeadline