package main

import (
	"context"
	"fmt"
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
		srv.Close()
		os.Exit(0)
	}()

	if err := srv.ListenAndServe(); err != nil{
		fmt.Println(err)
		cancel()
		os.Exit(1)
	}
}