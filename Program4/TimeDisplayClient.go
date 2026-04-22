/*
####################################################################################
#   Author: Jared Newton                                                           #
#   Date: 03/04/2026                                                               #
#   Description: A program written in go that connects to a chat server            #
#       and extracts a covert message based on the delay between characters sent   #
####################################################################################
*/
package main

import (
	"fmt"
	"net"
	"time"
)

const (
	address   = "138.47.99.21" // address of chat server
	port      = "31337"        // port of chat server
	threshold = 53.75         // midpoint between the two observed delays
)

func main() {
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Println("Failed to connect to server")
		return
	}

	defer conn.Close()
	fmt.Println("[Connected to chat server]")

	buffer := make([]byte, 1)
	var bits []int
	var message string
	var end = 0
	var charCount = 0 // tracks position so we can print delays in a readable way

	conn.Read(buffer)
	fmt.Printf("%c", buffer[0])
	previousTime := time.Now()

	for end != 3 {
		_, err := conn.Read(buffer)
		if err != nil {
			break
		}

		now := time.Now()
		delayMs := now.Sub(previousTime).Seconds() * 1000
		previousTime = now
		charCount++

		// Print character and its delay side by side
		fmt.Printf("%c[%.3fms] ", buffer[0], delayMs)

		// Print a newline every 8 characters to keep output readable
		if charCount%8 == 0 {
			fmt.Println()
		}

		if delayMs < threshold {
			bits = append(bits, 0)
		} else {
			bits = append(bits, 1)
		}

		if len(bits) == 8 {
			var value int
			for _, bit := range bits {
				value = (value << 1) | bit
			}
			char := rune(value)
			bits = bits[:0]
			message += string(char)

			if char == 'E' {
				end = 1
			} else if char == 'O' && end == 1 {
				end = 2
			} else if char == 'F' && end == 2 {
				end = 3
			} else {
				end = 0
			}
		}
	}

	fmt.Println()
	fmt.Println("[Closing connection to chat server]")
	fmt.Println("The covert message is:")
	fmt.Println(message[:len(message)-3])
}
