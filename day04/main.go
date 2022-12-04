package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type pair struct {
	firstStart  int
	firstEnd    int
	secondStart int
	secondEnd   int
}

func pairs() (pairs []pair) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	regex := regexp.MustCompile(`(\d+)-(\d+),(\d+)-(\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		match := regex.FindStringSubmatch(line)

		firstStart, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}
		firstEnd, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}
		secondStart, err := strconv.Atoi(match[3])
		if err != nil {
			log.Fatal(err)
		}
		secondEnd, err := strconv.Atoi(match[4])
		if err != nil {
			log.Fatal(err)
		}

		pair := pair{firstStart, firstEnd, secondStart, secondEnd}

		pairs = append(pairs, pair)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return pairs
}

func part1(pairs []pair) int {
	wasteful := 0
	for _, pair := range pairs {
		if (pair.firstStart <= pair.secondStart && pair.firstEnd >= pair.secondEnd) || (pair.secondStart <= pair.firstStart && pair.secondEnd >= pair.firstEnd) {
			wasteful++
		}
	}

	return wasteful
}

func part2(pairs []pair) int {
	wasteful := 0
	for _, pair := range pairs {
		if pair.firstStart <= pair.secondEnd && pair.firstEnd >= pair.secondStart {
			wasteful++
		}
	}

	return wasteful
}

func main() {
	pairs := pairs()
	fmt.Println(part1(pairs))
	fmt.Println(part2(pairs))
}
