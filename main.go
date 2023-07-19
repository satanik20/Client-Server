package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	// Enable Keepalives
	err = conn.SetKeepAlive(true)
	if err != nil {
		fmt.Printf("Unable to set keepalive - %s", err)
	}

	/*_, err = conn.Write([]byte("Heart Beat"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}*/

	j := 0
	duration := time.Duration(5000000000)
	msg := ""
	client_msg := "Start Heart Beat..."
	go func(conn net.Conn) {
		for range time.Tick(duration) {
			if j%4 == 0 {
				duration = time.Second * 20
				mem := &runtime.MemStats{}
				//cpu := runtime.NumCPU()
				runtime.ReadMemStats(mem)
				// Get CPU from pkg
				before, err := cpu.Get()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					return
				}
				time.Sleep(time.Duration(1) * time.Second)
				after, err := cpu.Get()
				if err != nil {
					log.Println(os.Stderr, "%s\n", err)
					return
				}
				total := float64(after.Total - before.Total)
				cpu_system := float64(after.User-before.User) / total * 100

				//msg = "CPU Utilization: " + strconv.Itoa(cpu) + " | Memory Utilisation: " + strconv.Itoa(int(mem.Alloc))
				msg = "CPU Utilization: " + strconv.Itoa(cpu_system) + " | Memory Utilisation: " + strconv.Itoa(int(mem.Alloc))
			} else {
				duration = time.Second * 5
				msg = "Heart Beat..."
				client_msg = "Heart Beat..."
			}
			fmt.Println(client_msg)
			_, err = conn.Write([]byte(msg))
			if err != nil {
				println("Write data failed:", err.Error())
				os.Exit(1)
			}

			j++
		}
	}(conn)

	time.Sleep(time.Second * 5)

	// CPU Utilization
	mem := &runtime.MemStats{}
	for {
		cpu := runtime.NumCPU()
		log.Println("Client CPU Utilization:", cpu)

		// Byte
		runtime.ReadMemStats(mem)
		log.Println("Memory Utilization:", mem.Alloc)

		time.Sleep(20 * time.Second)
		log.Println("----------")
	}
}
