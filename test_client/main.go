package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	urlf := flag.String("url", "ws://localhost:8080/", "url of the streaming server")

	flag.Parse()

	serverUrl := *urlf
	_, err := url.Parse(serverUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	conn, _, err := websocket.DefaultDialer.Dial(serverUrl, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer conn.Close()

	fmt.Println("connected")

	datachan := make(chan []byte)
	ctx, cancel := context.WithTimeout(context.Background(),time.Second * 3)
	go func() {
		<-ch
		close(ch)
		cancel()

	}()

	go func(ch chan []byte) {
		for{
			select{
			case <-ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					return
				}
				ch <- msg
			}
		}
	}(datachan)
	
	fmt.Println("starting to listen for messages")

	// enc := encoder.NewEncoder()
	cnt := 0
	for {
		select {
		case <-ctx.Done():
			close(datachan)
			fmt.Println("done: closing the connection")
			fmt.Println("received",cnt,"images")
			return
			
		case x := <-datachan:
			cnt+= 1
			// _,err := enc.Base64ToImage(x)
			if err != nil{
				fmt.Println(err)
				continue
			}

			f, err := os.Create(fmt.Sprintf("img_%d.jpg",cnt))
			if err != nil{
				fmt.Println(err)
				continue
			}
			f.Write(x)
			f.Close()
		}
	}

}
