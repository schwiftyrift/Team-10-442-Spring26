/*
* Author:	Kate Barron
* Date:		1st April, 2026
* Description: 	go code to log into an ftp server and read the file permissions from a given folder on the server
 */

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/secsy/goftp"
)

const (
	address  = "138.47.99.21" // the ip address of the ftp server
	port     = "21"           // the port (21 is the traditional port)
	username = "anonymous"    // what is the username
	password = ""             // what is the password
	path     = "/7"           // where are the files with the covert message
	METHOD   = 7              // either 7 or 10
)

func btoa(binary string) []byte {

	// variables
	var Result []byte
	var mode int
	var BitString string

	// if 7-bit method, split into 10-bit chunks
	if METHOD == 7 {
		mode = 10
	} else {
		// if 10-bit method, split into 7-bit chunks
		mode = 7
	}
	for i := 0; i < len(binary); i += mode {
		//get 7 or 10-bit chunk
		if METHOD == 7 {
			BitString = binary[i+2 : i+mode] // i + 2 cuts off initial "--" from chunk
		} else {
			BitString = binary[i : i+mode]
		}

		// Convert string chunk to ASCII (base 2 -> 64-bit)
		intValue, err := strconv.ParseInt(BitString, 2, 64)
		if err != nil {
			log.Fatal("Error convering binary to ascii: ", err)
		}

		// Append chunk to result
		Result = append(Result, byte(intValue))
	}
	return Result
}

func main() {
	config := goftp.Config{
		User:     username,
		Password: password,
		//ActiveTransfers: true,	// PASSIVE vs ACTIVE
	}

	// connect to FTP Server
	client, err := goftp.DialConfig(config, address+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	// List all the files
	entries, err := client.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// for each entry, print out its permissions, followed by the file name
	encodedStr := ""
	for _, entry := range entries {
		//fmt.Printf("%s\t%s\n", entry.Mode().String(), entry.Name())
		if METHOD == 7 {
			// Don't use entries that start with anything other than ---
			if strings.HasPrefix(entry.Mode().String(), "---") == true {
				encodedStr += entry.Mode().String()
			}
		} else if METHOD == 10 {
			encodedStr += entry.Mode().String()
		}

	}

	encodedBin := ""
	for _, letter := range encodedStr {
		// convert '-' to 0 and permissions to '1'
		if (letter == 'd') || (letter == 'r') || (letter == 'w') || (letter == 'x') {
			encodedBin += "1"
		} else {
			encodedBin += "0"
		}
	}
	//fmt.Println(encodedBin)

	fmt.Printf(string(btoa(encodedBin)) + "\n")

}
