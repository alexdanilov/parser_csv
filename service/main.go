package main

import (
	"net"
	"log"
	"bufio"
	"flag"
	"encoding/json"
	item "parser_service/item"
	"fmt"
	"sync"
)

type Database map[string]item.Item


var (
	data Database
	tcpAddress = flag.String("addr", "127.0.0.1:5000", "logger IP:port address")
)


func printDatabase() {
	fmt.Println("Database is:")

	for _, rec := range data {
		fmt.Println(rec)
	}
}


// Serve gets new connection and run loop to get messages from connection and puts to a queue
func serve(conn net.Conn) {
	var mu sync.RWMutex
	defer conn.Close()

	client := conn.RemoteAddr().String()
	log.Println("New client connected:", client)

	for {
		// listen for message ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Close connection for", client)
			go printDatabase()
			break
		}

		item := item.Item{}
		json.Unmarshal([]byte(message), &item)

		log.Println("Received new msg", message)
		if item.Id == "" {
			continue
		}
		log.Println("Received new item", item)

		mu.Lock()
		data[item.Id] = item
		mu.Unlock()
	}
}


func main() {
	flag.Parse()
	log.Println("Staring server...")

	data = make(Database)

	// listen port
	ln, _ := net.Listen("tcp", *tcpAddress)

	// run loop to accept clients connections
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(c)
	}
}
