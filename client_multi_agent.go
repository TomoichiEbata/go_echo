// go run client_multi_agent.go

// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")


var interrupt = make(chan os.Signal, 1) // Go のシグナル通知は、チャネルに os.Signal 値を送信することで行います。
			  	 		   			  	 // これらの通知を受信するためのチャネル (と、プログラムが終了できること
										 // を通知するためのチャネル) を作ります。
var	u = url.URL{Scheme: "ws", Host: *addr, Path: "/echo"} // JSON型の配列記述方法？

func main() {
	flag.Parse()     //引数の読み込み argv, argcと同じ
	log.SetFlags(0)  // ログの出力で時間の情報、この時点で0秒にセット

//	interrupt := make(chan os.Signal, 1) // Go のシグナル通知は、チャネルに os.Signal 値を送信することで行います。
//			  	 		   			  	 // これらの通知を受信するためのチャネル (と、プログラムが終了できること
//										 // を通知するためのチャネル) を作ります。

	signal.Notify(interrupt, os.Interrupt) // 指定されたシグナル通知を受信するために、 与えられたチャネルを登録
							 			   // します。

	log.Printf("connecting to %s", u.String()) // ここは単にプリントしているだけ(だろう)

	///// ここまでは共通 /////

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		

	go sub_main()
	time.Sleep(time.Millisecond * 333)  // 333ミリ秒		
	
	ticker := time.NewTicker(time.Second)  // 1秒おきに通知 (sleepと同じ)		

	for {
		select{
    	   case <-ticker.C: // tickerのチャネルはデフォルトで付いているらしい
		   case <-interrupt:  // こっちは手動割り込みだな検知だな
		        return;
		}
     }
}

func sub_main(){
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil) // これがコネクションの実施
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close() // deferは、どこに書かれていようとも、関数から抜ける前に実行される

	done := make(chan struct{}) // 配列といってもいいし、並行処理用のキューといってもいい
                                // 値が入っていないとデッドロックする

	go func() {  // 受信用スレッドを立ち上げる(スレッドの中でスレッド立ち上げているが、大丈夫だろうか)
		defer close(done)
		for {
			_, message, err := c.ReadMessage() // このメソッドの返り値は3つで、最初の返り値は不要
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)  // 受信したら、そのメッセージを表示する
		}
	}()

	ticker := time.NewTicker(time.Second)  // 1秒おきに通知 (sleepと同じ)
	defer ticker.Stop() // このループを抜ける時に終了する

	for {            // 無限ループの宣言かな(C/C++ で言うとろの、whileとかdoとか)
		select {
		case <-done: // doneの中に何かが入っていたら、このルーチンはリターンして終了となる
			 		 // (でも何も入っていないところを見ると、func()ルーチンの消滅で、こっちが起動するんだろう)
			return
		case t := <-ticker.C: // tickerのチャネルはデフォルトで付いているらしい
			   	  			  // 時間が入ってくるまでロックされる	 
			   	  			  // この場合1秒単位でチャネルに時間が放り込まれるのでそこで動き出す。

			err := c.WriteMessage(websocket.TextMessage, []byte(t.String())) // サーバに(時刻の)メッセージを送付する
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:  // こっちは手動割り込みだな検知だな
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.

			// close メッセージを送信してから、サーバーが接続を閉じるのを
			// （タイムアウトして）待つことで、接続をきれいに閉じます

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
				case <-done: // これも上記の"done"の説明と同じでいいかな
				case <-time.After(time.Second): // このメソッド凄い time.Secondより未来ならひっかかる
			}
			return
		}
	}
}

