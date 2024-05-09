package main

import (
	"io"
	"log"
	"os"
)

func main() {
	// ファイル作成
	file, err := os.Create("log.txt")
	if err != nil {
		// errの文字列を出力してプログラムを強制終了
		log.Fatalln(err)
	}
	flags := log.Lshortfile // log.Lshortfileでファイル名と行数を与える
	warnLogger := log.New(io.MultiWriter(file, os.Stderr), "WARN: ", flags)
	errorLogger := log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", flags)

	warnLogger.Println("warning A") // WARN: main.go:20: warning A

	errorLogger.Fatalln("critical error") // ERROR: main.go:22: critical error
}
