package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type instruction struct {
	count int
	from  int
	to    int
}

func parse() string {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		return scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ""
}

func findMarker(signal string, length int) int {
	for start := 0; start < len(signal)-length+1; start++ {
		marker := true
		for a := start; a < start+length-1; a++ {
			for b := a + 1; b < start+length; b++ {
				if signal[a] == signal[b] {
					marker = false
				}
			}
		}

		if marker {
			return start + length
		}
	}

	return -1
}

func part1(signal string) int {
	return findMarker(signal, 4)
}

func part2(signal string) int {
	return findMarker(signal, 14)
}

func main() {
	signal := parse()
	fmt.Println(part1(signal))
	fmt.Println(part2(signal))
}
