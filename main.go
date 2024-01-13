package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"screen_stream/screenmgr"
	"screen_stream/server"
	cfg "screen_stream/util/config"

	"syscall"
)

func main() {
	config, err := cfg.Load(".")

	if err != nil {
		log.Fatal("can't load config from the path .",)
	}

	log := log.New(os.Stdout,"[STREAM SERVER] ", 0)
	disp := screenmgr.NewDisplay(0)
	srv := server.New(config, log, disp)
	
	ch := make(chan os.Signal, 1)
	
	signal.Notify(ch, syscall.SIGINT)

	go func(){
		<-ch
		srv.Stop()
		os.Exit(0)
	}()

	http.HandleFunc("/screen",srv.SpawnNewScreenStream())
	http.HandleFunc("/events", srv.SpawnNewEventsHandler())
	
	if err := http.ListenAndServe(":8080",nil); err != nil{
		fmt.Println(err)
		srv.Stop()
		os.Exit(1)
	}
}