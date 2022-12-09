package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type move struct {
	direction string
	distance  int
}

type position struct {
	x int
	y int
}

func parse() []move {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var moves []move

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		distance, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal(err)
		}

		move := move{parts[0], distance}
		moves = append(moves, move)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return moves
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func moveHead(head position, direction string) position {
	switch direction {
	case "L":
		return position{head.x - 1, head.y}
	case "R":
		return position{head.x + 1, head.y}
	case "U":
		return position{head.x, head.y - 1}
	case "D":
		return position{head.x, head.y + 1}
	}

	return position{}
}

func moveTail(head position, tail position) position {
	dx := absDiffInt(head.x, tail.x)
	dy := absDiffInt(head.y, tail.y)

	if dx <= 1 && dy <= 1 {
		return tail
	}

	if head.x == tail.x {
		return position{tail.x, (head.y + tail.y) / 2}
	}

	if head.y == tail.y {
		return position{(head.x + tail.x) / 2, tail.y}
	}

	if head.x < tail.x {
		if head.y < tail.y {
			return position{tail.x - 1, tail.y - 1}
		} else {
			return position{tail.x - 1, tail.y + 1}
		}
	} else {
		if head.y < tail.y {
			return position{tail.x + 1, tail.y - 1}
		} else {
			return position{tail.x + 1, tail.y + 1}
		}
	}
}

func part1(moves []move) int {
	head := position{0, 0}
	tail := position{0, 0}

	positions := map[position]bool{tail: true}
	for _, move := range moves {
		for i := 0; i < move.distance; i++ {
			head = moveHead(head, move.direction)
			tail = moveTail(head, tail)
			positions[tail] = true
		}
	}

	return len(positions)
}

func part2(moves []move) int {
	var knots []position
	for i := 0; i < 10; i++ {
		knots = append(knots, position{0, 0})
	}

	positions := map[position]bool{knots[9]: true}
	for _, move := range moves {
		for i := 0; i < move.distance; i++ {
			knots[0] = moveHead(knots[0], move.direction)
			for j := 1; j < 10; j++ {
				knots[j] = moveTail(knots[j-1], knots[j])
			}
			positions[knots[9]] = true
		}
	}

	return len(positions)
}

func main() {
	moves := parse()
	fmt.Println(part1(moves))
	fmt.Println(part2(moves))
}
