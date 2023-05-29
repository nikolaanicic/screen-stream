package benchmarks

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"screen_stream/server"
	"testing"
	"time"
)




func BenchmarkServer(b *testing.B) {

	ctx, cancel := context.WithTimeout(context.Background(),time.Minute)
	srv := server.New(ctx,cancel,log.Default())
	
	mux := http.NewServeMux()
	mux.HandleFunc("/",srv.SpawnNewStream())

	httpsrv := &http.Server{Addr: ":8080",Handler: mux}
	
	
	go func(){
		if err := httpsrv.ListenAndServe(); err != nil{
			fmt.Println(err)
			cancel()
			return
		}
	}()


	timer := time.NewTimer(time.Minute)
	<-timer.C
	srv.Stop()
	httpsrv.Shutdown(ctx)

}