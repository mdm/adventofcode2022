package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type instruction struct {
	count int
	from  int
	to    int
}

func parse() (stacks [][]rune, instructions []instruction) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	instructionRegex := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)

	scanner := bufio.NewScanner(file)
	phase1 := true
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			phase1 = false
			continue
		}

		if phase1 {
			for i := 0; i*4+1 < len(line); i++ {
				if len(stacks) <= i {
					var stack []rune
					stacks = append(stacks, stack)
				}

				mark := []rune(line)[i*4+1]

				if mark >= 'A' && mark <= 'Z' {
					stacks[i] = append(stacks[i], mark)
				}
			}
		} else {
			match := instructionRegex.FindStringSubmatch(line)

			count, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}

			from, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal(err)
			}

			to, err := strconv.Atoi(match[3])
			if err != nil {
				log.Fatal(err)
			}

			instruction := instruction{count, from - 1, to - 1}
			instructions = append(instructions, instruction)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return stacks, instructions
}

func part1(stacks [][]rune, instructions []instruction) string {
	for _, instruction := range instructions {
		for i := 0; i < instruction.count; i++ {
			crate := stacks[instruction.from][0]
			stacks[instruction.from] = stacks[instruction.from][1:]
			stacks[instruction.to] = append([]rune{crate}, stacks[instruction.to]...)
		}
	}

	var message []rune
	for _, stack := range stacks {
		message = append(message, stack[0])
	}

	return string(message)
}

func part2(stacks [][]rune, instructions []instruction) string {
	for _, instruction := range instructions {
		var crates []rune
		for i := 0; i < instruction.count; i++ {
			crates = append(crates, stacks[instruction.from][0])
			stacks[instruction.from] = stacks[instruction.from][1:]
		}
		stacks[instruction.to] = append(crates, stacks[instruction.to]...)
	}

	var message []rune
	for _, stack := range stacks {
		message = append(message, stack[0])
	}

	return string(message)
}

func main() {
	stacks, instructions := parse()
	fmt.Println(part1(stacks, instructions))
	stacks, instructions = parse()
	fmt.Println(part2(stacks, instructions))
}
