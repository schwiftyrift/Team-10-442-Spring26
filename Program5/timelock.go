/*
* Team: 10 -- Titan
* Date: 4/27/26
 */

package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	DEBUG = false
)

func main() {

	var year, month, day, hour, minute, second int

	_, err := fmt.Fscanf(os.Stdin, "%d %d %d %d %d %d",
		&year, &month, &day,
		&hour, &minute, &second)
	if err != nil {
		log.Fatal("Invalid input format:", err)
	}

	// set local time zone to chicago timezone, IDK apparently my local timezone was wrong on my VM so this is a fix for that
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatal(err)
	}

	// convert given time to time.Date
	epochTime := time.Date(
		year,
		time.Month(month),
		day,
		hour,
		minute,
		second,
		0,
		loc,
	).UTC() // convert to UTC (doesn't use daylight savings)

	var currentTime time.Time

	// if in Debug mode, use a set date placed here
	if DEBUG {
		currentTime = time.Date(2015, 05, 15, 14, 00, 00, 0, loc).UTC() // UTC ignores time zones entirely
		fmt.Println("DEBUG: using fixed current time")
	} else {
		currentTime = time.Now().UTC() // get current time and convert to UTC
	}

	// get the difference between current and epoch time
	elapsedTime := int64(currentTime.Sub(epochTime).Seconds())
	if DEBUG {
		fmt.Println("Elapsed time (seconds): ", elapsedTime)
	}

	// Floor to 60-second interval
	intervalTime := elapsedTime - (elapsedTime % 60)

	// md5 hash the elapsed time twice
	input := strconv.FormatInt(intervalTime, 10)
	h1 := md5.Sum([]byte(input))
	h2 := md5.Sum([]byte(hex.EncodeToString(h1[:])))
	finalHash := hex.EncodeToString(h2[:])

	if DEBUG {
		fmt.Println("Hash:", finalHash)
	}

	// Get first two letters from left (a - f)
	letters := ""
	for _, c := range finalHash {
		if c >= 'a' && c <= 'f' {
			letters += string(c)
			// only get first 2 letters
			if len(letters) == 2 {
				break
			}
		}
	}

	// get first two numbers from right (0 - 9)
	numbers := ""
	for i := len(finalHash) - 1; i >= 0; i-- {
		c := rune(finalHash[i])
		if c >= '0' && c <= '9' {
			numbers += string(c)
			if len(numbers) == 2 {
				break
			}
		}
	}

	fmt.Println(letters + numbers) // print final code

}
