package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"screen_stream/server"
	cfg "screen_stream/util/config"

	"syscall"
)

func main() {
	config, err := cfg.Load(".")

	if err != nil {
		log.Fatal("can't load config from the path .",)
	}

	srv := server.New(config,log.Default())
	
	ch := make(chan os.Signal, 1)
	
	signal.Notify(ch, syscall.SIGINT)

	go func(){
		<-ch
		srv.Stop()
		os.Exit(0)
	}()

	http.HandleFunc("/",srv.SpawnNewStream())
	
	if err := http.ListenAndServe(":8080",nil); err != nil{
		fmt.Println(err)
		srv.Stop()
		os.Exit(1)
	}
}