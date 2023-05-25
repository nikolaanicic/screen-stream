package main

import (
	"log"
	"screen_stream/test_client/app"
)

func main() {
	url := "ws://localhost:8080/"
	
	app := app.NewApp("streamer",log.Default())

	if err := app.Start(url); err != nil{
		log.Fatal(err)
	}
}
