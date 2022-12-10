package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parse() []string {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var instructions []string

	for scanner.Scan() {
		line := scanner.Text()

		instructions = append(instructions, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return instructions
}

func part1(instructions []string) int {
	sum := 0

	x := 1
	cycle := 0

	for _, instruction := range instructions {
		var remaining int
		if instruction == "noop" {
			remaining = 1
		} else {
			remaining = 2
		}

		for i := 0; i < remaining; i++ {
			cycle++

			if (cycle-20)%40 == 0 {
				sum += cycle * x
			}
		}

		if instruction != "noop" {
			operand, err := strconv.Atoi(strings.Split(instruction, " ")[1])
			if err != nil {
				log.Fatal(err)
			}

			x += operand
		}
	}

	return sum
}

func part2(instructions []string) {
	x := 1
	cycle := 0

	for _, instruction := range instructions {
		var remaining int
		if instruction == "noop" {
			remaining = 1
		} else {
			remaining = 2
		}

		for i := 0; i < remaining; i++ {
			cycle++

			if (cycle-1)%40 >= x-1 && (cycle-1)%40 <= x+1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}

			if cycle%40 == 0 {
				fmt.Println()
			}
		}

		if instruction != "noop" {
			operand, err := strconv.Atoi(strings.Split(instruction, " ")[1])
			if err != nil {
				log.Fatal(err)
			}

			x += operand
		}
	}
}

func main() {
	instructions := parse()
	fmt.Println(part1(instructions))
	part2(instructions)
}
