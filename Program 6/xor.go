/*
####################################################################################
#   Author: Jared Newton                                                           #
#   Date: 04/28/2026                                                               #
#   Description: A program written in go that takes a file input from stdin        #
#       and performs an XOR with a key file                                        #
####################################################################################
*/

package main

import (
	"fmt"
	"io"
	"os"
)

const keyPath = "key"

func xor_function(input, key []byte) []byte {
	var result []byte
	for i, j := range input {
		// xor input with key, wrap key around if it is smaller than the input
		result = append(result, j^key[i%len(key)])
	}
	return result
}

func main() {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin: %v", err)
		os.Exit(1)
	}

	fileBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin: %v", err)
		os.Exit(1)
	}

	os.Stdout.Write(xor_function(fileBytes, keyBytes))
}
