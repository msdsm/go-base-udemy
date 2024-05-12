package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	file, err := os.Create("log.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	errorLogger := log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.LstdFlags)
	ctx, cancel := context.WithTimeout(context.Background(), 5100*time.Millisecond)
	defer cancel()
	const wdtTimeout = 800 * time.Millisecond   // watchdog timerのタイムアウト
	const beatInterval = 500 * time.Millisecond // heartbeatの周期
	heartbeat, v := task(ctx, beatInterval)
loop:
	for {
		select {
		case _, ok := <-heartbeat: // heartbeatから送信があったとき
			if !ok { // heartbeat channelがcloseしたとき
				break loop
			}
			fmt.Println("beat pulse ")
		case r, ok := <-v: // valueの書き込み
			if !ok { // value channelがclose
				break loop
			}
			t := strings.Split(r.String(), "m=")
			fmt.Printf("value: %v [s]\n", t[1])
		case <-time.After(wdtTimeout): // watchdog timerのタイムアウト(タイムアウト時間の間に上のcase文がどれも実行されていないとここに到達)
			errorLogger.Println("doTask goroutine's heartbeat stopped")
			break loop
		}
	}
}

func task(
	ctx context.Context,
	beatInterval time.Duration,
) (<-chan struct{}, <-chan time.Time) { // 返り値はheartbeat, value(通知専用channel, 時刻型の読み取り専用channel)
	heartbeat := make(chan struct{}) // heartbeat専用のデータなしchannel
	out := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(out)
		pulse := time.NewTicker(beatInterval)
		task := time.NewTicker(2 * beatInterval)
		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}: // heartbeatが受信を開始している場合
			default:
			}
		}
		sendValue := func(t time.Time) {
			for {
				select {
				case <-ctx.Done(): // contextのタイムアウトに到達したとき
					return
				case <-pulse.C: // heartbeatが送信されているとき
					sendPulse()
				case out <- t:
					return
				}
			}
		}
		// var i int
		for {
			select {
			case <-ctx.Done():
				return
			case <-pulse.C: // heartbeatのchannelに書き込みがあった時にsendPulseを呼び出す
				// if i == 3 {
				// 	time.Sleep(1000 * time.Millisecond)
				// }
				sendPulse()
				// i++
			case t := <-task.C: // taskに書き込みがあった時
				sendValue(t)
			}
		}
	}()
	return heartbeat, out
}
