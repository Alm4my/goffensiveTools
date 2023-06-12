package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
)

const DefaultWorkerNbr = 10000
const DefaultIPAddress = "127.0.0.1"
const DefaultPort = "-1"

func worker(ports chan int, results chan int) {
	for port := range ports {
		adresse := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.Dial("tcp", adresse)
		if err != nil {
			results <- 0
			continue
		}
		err = conn.Close()
		if err != nil {
			fmt.Println("There was an error closing the address", adresse)
		}
		results <- port
	}
}

func main() {

	// flags
	var ipAddress string
	var port string
	var workers int

	// flags declaration
	flag.StringVar(&ipAddress, "i", DefaultIPAddress, "Specify an IP address. The default is 127.0.0.1.")
	flag.StringVar(&port, "p", DefaultPort, "Specify a port. The default is -1 to scan all ports.")
	flag.IntVar(&workers, "w", DefaultWorkerNbr, "Specify the number of workers. The default is 10,000")

	// calling flags
	flag.Parsed()

	// create two channels and a slice for communication
	ports := make(chan int, workers)
	results := make(chan int)
	var openPorts []int

	// number of times we open a subroutine to buffer
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	// send ports to ports channel x times
	portRange := 0
	if port == "-1" {
		portRange = 65_535
	}
	go func() {
		for i := 1; i <= portRange; i++ {
			ports <- i
		}
	}()

	// gather results from results channel
	for i := 0; i < portRange; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Println("Le port", port, "de l'adresse", DefaultIPAddress, "est ouvert.")
	}
}
