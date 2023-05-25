package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image/jpeg"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
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
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ch
		cancel()

	}()

	go func(ch chan []byte) {
		for {
			select {
			case <-ctx.Done():
				close(datachan)
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

	a := app.New()
	w := a.NewWindow("My new window")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(840, 680))

	w.SetOnClosed(func() {
		cancel()
	})

	cnt := 0
	go func() {
		samplePeriod := float64(1000) / float64(30)
		ticker := time.NewTicker(time.Millisecond * time.Duration(samplePeriod))

		for {
			select {
			case <-ctx.Done():
				close(datachan)
				fmt.Println("done: closing the connection")
				fmt.Println("received", cnt, "images")
				
				return

			case x := <-datachan:
				cnt++
				img, err := jpeg.Decode(bytes.NewReader(x))
				if err != nil {
					fmt.Println(err)
					continue
				}

				fimg := canvas.NewImageFromImage(img)
				fimg.FillMode = canvas.ImageFillContain
				<-ticker.C
				w.SetContent(fimg)

			}
		}
	}()


	w.ShowAndRun()
}
