package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type strategy struct {
	challenge int
	response  int
}

func strategies() (strategies []strategy) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")

		challenge := int(parts[0][0]) - int('A')
		response := int(parts[1][0]) - int('X')

		strategy := strategy{challenge, response}
		strategies = append(strategies, strategy)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return strategies
}

func part1(strategies []strategy) int {
	score := 0
	for _, strategy := range strategies {
		score += strategy.response + 1

		if (strategy.challenge+1)%3 == strategy.response {
			score += 6
		}

		if strategy.challenge == strategy.response {
			score += 3
		}
	}

	return score
}

func part2(strategies []strategy) int {
	score := 0
	for _, strategy := range strategies {
		var choice int

		switch strategy.response {
		case 0:
			choice = (strategy.challenge + 2) % 3
		case 1:
			choice = strategy.challenge
			score += 3
		case 2:
			choice = (strategy.challenge + 1) % 3
			score += 6
		}

		score += choice + 1
	}

	return score
}

func main() {
	strategies := strategies()
	fmt.Println(part1(strategies))
	fmt.Println(part2(strategies))
}
