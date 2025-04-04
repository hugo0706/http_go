package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hugo0706/httpfromtcp/internal/request"
)

func main(){
  l, err := net.Listen("tcp", "localhost:42069")
  if err != nil {
  	log.Fatalf("Cannot listen, error: %s", err.Error())
  }

  for {
  	con, err := l.Accept()
   	if err != nil {
    	fmt.Printf("Error: %s", err.Error())	
    	break
    }
   	fmt.Print("Connection stablished!\n")
    
    req, err := request.RequestFromReader(con)
    if err != nil {
   		fmt.Printf("Error: %s", err.Error())	
    }
    
    fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", 
    	req.RequestLine.Method,
    	req.RequestLine.RequestTarget,
    	req.RequestLine.HttpVersion)
		defer fmt.Println("Connection closed!")
  }
  
}

