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
	address = "138.47.99.21" // address of chat server
	port    = "31337"      // port of chat server
	threshold = 53.75     // midpoint between observed delay times
)

func main() {
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Println("Failed to connect to server")
	}

	defer conn.Close()
	fmt.Println("[Connected to chat server]")
	buffer := make([]byte, 1) // create buffer to read bytes from the chat server
	var bits []int            // array to store detected bits
	var message string = ""   // final covert message
	var end = 0               // keeps track of when to end the connection

	conn.Read(buffer)            // read in first character
	fmt.Print(string(buffer[0])) // print read in character
	previousTime := time.Now()   // set he time to when first character is read in

	for end != 3 {
		_, err := conn.Read(buffer)
		if err != nil {
			break
		}

		now := time.Now()
		delay := now.Sub(previousTime).Seconds() * 1000 // calculate delay between characters
		previousTime = now
		fmt.Print(string(buffer[0]))
		var char rune

		// determine if delay corresponds to 1 or 0
		if delay < threshold {
			bits = append(bits, 0)
		} else {
			bits = append(bits, 1)
		}

		// when 8 bits are read in, convert it to ASCII and add to message
		if len(bits) == 8 {
			var value int
			for _, bit := range bits {
				value = (value << 1) | bit
			}
			char = rune(value)
			bits = bits[:0]
			message += string(char)

			// check to see if "EOF" is found
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
	conn.Close()
	fmt.Println()
	fmt.Println("[Closing connection to chat server]")
	fmt.Println("The covert message is:")
	fmt.Println(message[:len(message)-3]) // print message without "EOF"

}
