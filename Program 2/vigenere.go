package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Function for encrypting the phrase
func encrypt(phrase []rune, key []rune) []rune {
	var result []rune
	// loop through each letter in phrase
	for i := 0; i < len(phrase); i++ {

		erune := phrase[i]
		if unicode.IsLetter(erune) { // only shift if the item is a letter (keep numbers/symbols the same)

			var base rune

			// check if value is uppercase or lowercase
			if unicode.IsUpper(erune) {
				base = 'A'
			} else {
				base = 'a'
			}

			shift := unicode.ToLower(key[i%len(key)]) - 'a'

			// apply shift
			erune = (erune-base+shift)%26 + base
		}
		result = append(result, erune)
	}
	return result
}

// function for decrypting a phrase
func decrypt(phrase []rune, key []rune) []rune {
	var result []rune

	// iterate through each letter/symbol in phrase
	for i := 0; i < len(phrase); i++ {
		drune := phrase[i]
		if unicode.IsLetter(drune) {
			var base rune

			if unicode.IsUpper(drune) { //check if letter is uppercase or lowercase
				base = 'A'
			} else {
				base = 'a'
			}

			// apply shift
			shift := unicode.ToLower(key[i%len(key)]) - 'a'
			drune = (drune-base-shift+26)%26 + base
		}
		result = append(result, drune)
	}
	return result
}

func main() {
	// mode and key are given with the command as args
	mode := os.Args[1]
	keyStr := os.Args[2]

	key := []rune(keyStr)

	// read for phrase and remove whitespace and newlines at the end
	reader := bufio.NewReader(os.Stdin)

	phraseStr, _ := reader.ReadString('\n')
	phraseStr = strings.TrimSpace(phraseStr)

	phrase := []rune(phraseStr)

	if mode == "-d" {
		fmt.Println(string(decrypt(phrase, key)))
	}
	if mode == "-e" {
		fmt.Println(string(encrypt(phrase, key)))
	}
}
