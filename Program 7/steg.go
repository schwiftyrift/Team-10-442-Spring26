/*
####################################################################################
#   Author: Team 10 - Titan                                                        #
#   Date: 05/06/2026                                                               #
#   Description: A program written in go that uses stegenography algorithms        #
#       to encode and decode files into and from a wrapper file                    #
####################################################################################
*/

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

const read_direction = "left"

var sentinel = [6]byte{0x0, 0xff, 0x0, 0x0, 0xff, 0x0}

func byte_method_encode(wrapperFile []byte, hiddenFile []byte, offset int, sentinel []byte, direction string, interval int) []byte {
	if interval == 1 {
		interval = int(math.Floor((float64(len(wrapperFile)-offset) / (float64(len(hiddenFile) + len(sentinel))))))
	}

	if direction == "left" {
		for i := 0; i < len(hiddenFile); i++ {
			wrapperFile[offset] = hiddenFile[i] //
			offset += interval
		}

		for i := 0; i < len(sentinel); i++ {
			wrapperFile[offset] = sentinel[i]
			offset += interval
		}
	} else if direction == "right" {
		offset = len(wrapperFile) - 1

		for i := 0; i < len(hiddenFile); i++ {
			wrapperFile[offset] = hiddenFile[i]
			offset -= interval
		}

		for i := 0; i < len(sentinel); i++ {
			wrapperFile[offset] = sentinel[i]
			offset -= interval
		}

	} else {
		fmt.Fprintf(os.Stderr, "Direction must be 'left' or 'right'")
		os.Exit(1)
	}

	return wrapperFile
}

func bit_method_encode(wrapperFile []byte, hiddenFile []byte, offset int, sentinel []byte, direction string) []byte {
	interval := 8
	if direction == "left" {
		for i := 0; i < len(hiddenFile); i++ {
			for j := 0; j < 8; j++ {
				wrapperFile[offset] &= 0xFE
				wrapperFile[offset] |= (hiddenFile[i] & 0x80) >> 7
				hiddenFile[i] <<= 1
				offset += int(interval)
			}
		}
		for i := 0; i < len(sentinel); i++ {
			for j := 0; j < 8; j++ {
				wrapperFile[offset] &= 0xFE
				wrapperFile[offset] |= (sentinel[i] & 0x80) >> 7
				sentinel[i] <<= 1
				offset += int(interval)
			}
		}
	} else if direction == "right" {
		for i := 0; i < len(hiddenFile); i++ {
			for j := 0; j < 8; j++ {
				wrapperFile[offset] &= 0xFE
				wrapperFile[offset] |= hiddenFile[i] & 0x01
				hiddenFile[i] >>= 1
				offset += int(interval)
			}
		}
		for i := 0; i < len(sentinel); i++ {
			for j := 0; j < 8; j++ {
				wrapperFile[offset] &= 0xFE
				wrapperFile[offset] |= (sentinel[i] & 0x01)
				sentinel[i] >>= 1
				offset += int(interval)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Direction must be 'left' or 'right'")
		os.Exit(1)
	}
	return wrapperFile
}

func byte_method_decode(wrapperFile []byte, offset int, sentinel []byte, interval int) []byte {
	var decoded []byte
	matched := 0

	for offset < len(wrapperFile) {
		b := wrapperFile[offset]

		if b == sentinel[matched] {
			matched++
			if matched == len(sentinel) {
				break // sentinel found, stop extracting bytes
			}
		} else {
			if matched > 0 {
				decoded = append(decoded, sentinel[:matched]...)
				matched = 0
			}
			decoded = append(decoded, b)
		}
		offset += interval
	}

	return decoded
}

func bit_method_decode(wrapperFile []byte, offset int, sentinel []byte, interval int) []byte {
	var decoded []byte
	matched := 0

	for offset < len(wrapperFile) {
		b := byte(0)

		for j := 0; j < 8; j++ {
			if offset >= len(wrapperFile) {
				break
			}
			b <<= 1
			b |= (wrapperFile[offset] & 0x01)
			offset += interval
		}

		if b == sentinel[matched] {
			matched++
			if matched == len(sentinel) {
				break
			}
		} else {
			if matched > 0 {
				decoded = append(decoded, sentinel[:matched]...)
				matched = 0
			}
			decoded = append(decoded, b)
		}
	}
	return decoded
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -(s|r) -(b|B) -o <val> [-i <val>] -w <val> [-h <val>]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  -s\t\tstore mode")
		fmt.Fprintln(os.Stderr, "  -r\t\tretrieve mode")
		fmt.Fprintln(os.Stderr, "  -b\t\tbit mode")
		fmt.Fprintln(os.Stderr, "  -B\t\tbyte mode")
		fmt.Fprintln(os.Stderr, "  -o <val>\tset offset (default 0)")
		fmt.Fprintln(os.Stderr, "  -i <val>\tset interval (default 1)")
		fmt.Fprintln(os.Stderr, "  -w <val>\twrapper file")
		fmt.Fprintln(os.Stderr, "  -h <val>\thidden file")
	}

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	store := flag.Bool("s", false, "store mode")
	retrieve := flag.Bool("r", false, "retrieve mode")
	bitMode := flag.Bool("b", false, "bit mode")
	byteMode := flag.Bool("B", false, "byte mode")
	offset := flag.Int("o", 0, "offset")
	interval := flag.Int("i", 1, "interval")
	wrapper := flag.String("w", "", "wrapper file")
	hidden := flag.String("h", "", "hidden file")

	flag.Parse()

	if *wrapper == "" {
		fmt.Fprintln(os.Stderr, "wrapper file is required (-w)")
		os.Exit(1)
	}
	if !*store && !*retrieve {
		fmt.Fprintln(os.Stderr, "must specify -s or -r")
		os.Exit(1)
	}
	if !*bitMode && !*byteMode {
		fmt.Fprintln(os.Stderr, "must specify -b or -B")
		os.Exit(1)
	}

	wrapperFile, err := os.ReadFile(*wrapper)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}

	var hiddenFile []byte
	if *store {
		if *hidden == "" {
			fmt.Fprintln(os.Stderr, "hidden file is required for store mode (-h)")
			os.Exit(1)
		}
		hiddenFile, err = os.ReadFile(*hidden)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
			os.Exit(1)
		}
	}

	if *retrieve && *bitMode {
		result := bit_method_decode(wrapperFile, *offset, sentinel[:], *interval)
		os.Stdout.Write(result)
	} else if *retrieve && *byteMode {
		result := byte_method_decode(wrapperFile, *offset, sentinel[:], *interval)
		os.Stdout.Write(result)
	} else if *store && *bitMode {
		result := bit_method_encode(wrapperFile, hiddenFile, *offset, sentinel[:], read_direction)
		os.Stdout.Write(result)
	} else if *store && *byteMode {
		result := byte_method_encode(wrapperFile, hiddenFile, *offset, sentinel[:], read_direction, *interval)
		os.Stdout.Write(result)
	}

}
