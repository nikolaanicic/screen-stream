package benchmarks

import (
	"context"
	"log"
	"net/http"
	"screen_stream/screenmgr"
	"screen_stream/server"
	cfg "screen_stream/util/config"

	"testing"
	"time"
)




func BenchmarkServer(b *testing.B) {

	config, err := cfg.Load("../..")

	if err != nil{
		log.Fatal("can't load config from the path")
	}

	disp := screenmgr.NewDisplay(0)
	srv := server.New(config,log.Default(), disp)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/",srv.SpawnNewScreenStream())

	httpsrv := &http.Server{Addr: ":8080",Handler: mux}
	
	
	go func(){
		if err := httpsrv.ListenAndServe(); err != nil{
			log.Println(err)
			return
		}
	}()


	timer := time.NewTimer(time.Minute)
	<-timer.C
	srv.Stop()
	httpsrv.Shutdown(context.Background())

}