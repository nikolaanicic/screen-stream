package benchmarks

import (
	"context"
	"log"
	"net/http"
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

	srv := server.New(config,log.Default())
	
	mux := http.NewServeMux()
	mux.HandleFunc("/",srv.SpawnNewStream())

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