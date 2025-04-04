package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main(){
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Cannot obtain address error: %s", err.Error())
		panic(err)
	}
	
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Cannot create connection error: %s", err.Error())
		panic(err)
	}
	defer conn.Close()
	
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				_, err = conn.Write([]byte(line))
				if err != nil {
					fmt.Printf("Error when writing: %s", err.Error())
				}
				break
			} else {
				fmt.Printf("Error when reading: %s", err.Error())
				break
			}
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Printf("Error when writing: %s", err.Error())
		}
		
	}
}