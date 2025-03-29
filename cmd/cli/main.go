package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	address := flag.String("address", "127.0.0.1:3223", "Server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Enter commands (SET key value, GET key, DEL key, or 'exit' to quit):")

	scanner := bufio.NewScanner(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Error sending command:", err)
			continue
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}
		fmt.Println(response)
	}
}
