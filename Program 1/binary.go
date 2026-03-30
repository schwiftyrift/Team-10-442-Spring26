package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func detectEncoding(input string) int {
	input = strings.TrimSpace(input)
	if len(input)%8 == 0 {
		return 8
	}
	if len(input)%7 == 0 {
		return 7
	}

	return -1
}
func parseBinary(input string) {
	input = strings.TrimSpace(input)
	message := ""
	encoding := detectEncoding(input)
	if encoding == 8 {
		for i := 0; i+8 <= len(input); i += 8 {
			chunk := input[i : i+8]
			temp, err := strconv.ParseInt(chunk, 2, 64)
			if err != nil {
				return
			}
			if temp == 8 {
				if len(message) > 0 {
					message = message[:len(message)-1]
				}
			} else {
				message += string(rune(temp))
			}
		}
		fmt.Println(message)
	} else {
		for i := 0; i+7 <= len(input); i += 7 {
			chunk := input[i : i+7]
			temp, err := strconv.ParseInt(chunk, 2, 64)
			if err != nil {
				return
			}
			if temp == 8 {
				if len(message) > 0 {
					message = message[:len(message)-1]
				}
			} else {
				message += string(rune(temp))
			}
		}
		fmt.Println(message)
	}

}

func main() {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading stdin:", err)
		return
	}
	input := string(content)
	parseBinary(input)
}
