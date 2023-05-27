package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"screen_stream/server"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	srv := server.NewServer(ctx,cancel)

	
	ch := make(chan os.Signal,1)
	signal.Notify(ch, syscall.SIGINT)

	go func(){
		<-ch
		srv.Stop()
		cancel()
		os.Exit(0)
	}()

	http.HandleFunc("/",srv.SpawnNewStream())
	
	if err := http.ListenAndServe(":8080",nil); err != nil{
		fmt.Println(err)
		cancel()
		os.Exit(1)
	}
}