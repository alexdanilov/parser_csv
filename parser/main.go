package main

import (
	"os"
	"encoding/csv"
	"bufio"
	"flag"
	"fmt"
	. "io"
	"time"
	"net"
	"log"
	"encoding/json"
	item "parser_service/item"
)


type itemQueue chan item.Item


var (
	bufferSize = flag.Int("buffer", 100, "a CSV buffer size")
	filePtr = flag.String("file", "data.csv", "a CSV file with data")
	phoneCode = flag.String("code", "+44", "Default phone code")
	tcpAddress = flag.String("addr", "127.0.0.1:5000", "logger IP:port address")
)


// parse given file to a struct and put it into queue
func parseFile(filename string, phoneCode string, q itemQueue) {
	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// open CSV reader
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// run reader loop
	for {
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		i := item.Item{record[0], record[1], record[2], record[3]}
		i.Normalize(phoneCode)
		q <- i
	}
}


// Returns new connection or wait until connect
func getConnection() (conn net.Conn) {
	var err error
	for {
		conn, err = net.Dial("tcp", *tcpAddress)
		if err != nil {
			log.Println("Cant connect to service:", err.Error())
			time.Sleep(time.Second)
			continue
		}
		break
	}

	log.Println("Connected to:", *tcpAddress)
	return conn
}


// send data from queue to an external service
func sendData(queue itemQueue) {
	var conn net.Conn

	defer func() {
		if conn == nil {
			return
		}
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	conn = getConnection()

	// run loop to read values from queue
	for {
		item := <-queue
		log.Println("Send new item", item)

		data, _ := json.Marshal(item)
		_, err := conn.Write(data)
		if err != nil {
			queue <- item
			conn = getConnection()
		}
		conn.Write([]byte("\n"))
	}
}


func main() {
	flag.Parse()

	requestsQueue := make(itemQueue, *bufferSize)
	go sendData(requestsQueue)

	// open and parse file
	parseFile(*filePtr, *phoneCode, requestsQueue)

	// wait for sending all data
	for {
		if len(requestsQueue) == 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
