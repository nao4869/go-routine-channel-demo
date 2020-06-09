package main

import (
	"fmt"
	"net/http"
	"time"
)

/* Go Routine
いつでもGo言語のプログラムを実行した際には、自動的に１つのGo Routineが作成される - プログラムのプロセスのようなこと

go キーワード - 非同期処理のようなもの？
	goキーワードを使用した際には、新たなthread Go Routineが作成される
	今回のケースでは、http.Get(link)のレスポンスを待っている間に、go checkLink(link)を並行処理で処理している

Go Scheduler
	１つのCPU Coreごとに起動している
	Go Routineは２つ以上同時に処理がされない
	Go Schedulerの役割は、どのGo Routineが現在起動しているか、次にどのRoutineを起動するかなどをモニタリングする

並行処理
	A、B、Cの処理をよーいどん、で3つ同時にスタートする
	複数のプロセス上で、複数のスレッドが立ち上がる
	一つのプロセス上で、複数のスレッドを切り替えながら処理をすることは並行処理という。
	私が初めて触った言語のJavaやGolang, Pythonなどで並列処理が可能
*/

func main() {
	// この書き方では、上から順に実行され、１つごとのリンクの状態を確認する際に遅延が起こる
	/*
		1. 最初のリンクをsliceから受け取る
		2. リクエストを送信
		3. レスポンスを待ち、出力
	*/
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://amazon.com",
		"http://golang.org",
	}

	// Channelを作成
	channel := make(chan string)

	for _, link := range links {
		// go キーワードは関数呼び出しの前のみで使用できる
		// main routineが終了した際に、child routineがあるかどうかなどは気にしない
		go checkLink(link, channel)
	}

	// 2. メッセージをListenし、受信次第printする
	// fmt.Println(<- channel)

	// 通常のforループを使用し、linksの数だけchannelで受信した値を出力
	// for i := 0; i < len(links); i++ {
	// 	fmt.Println(<- channel)
	// }
	
	for l := range channel {
		// 関数リテラル, Dart -> () {}
		// 新規Go Routineを開始するためには、go funcの形式にする
		go func(link string) {
			time.Sleep(3 * time.Second)
			go checkLink(link, channel)
		}(l)
	}
}

func checkLink(link string, channel chan string) {
	// この関数からレスポンスが来るのを待機
	// 1000個のステータスチェックをする際などには、膨大な時間がかかる
	// 解決策 - 並行処理を実装する
	_, error := http.Get(link)

	// エラーハンドリング
	if error != nil {
		fmt.Println(link + "might be down!")
		channel <- "Might be down" // 1. stringメッセージをchannelへ送信
		return
	}

	fmt.Println(link, "is up!")
	channel <- "It's up" // 1. stringメッセージをchannelへ送信
}
