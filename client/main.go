package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	urlf := flag.String("url","localhost:8080/","url of the streaming server")

	flag.Parse()

	serverUrl, err := url.Parse(*urlf)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	
}