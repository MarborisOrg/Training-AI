package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8081")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	message := "req=false,rate=0.1,hiddensNodes=50"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	buffer := make([]byte, 1024)
	recv_result, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving message:", err)
	} else {
		fmt.Println("Server Response:", string(buffer[:recv_result]))
	}
}
