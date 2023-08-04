package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type vec2 struct {
	x int
	y int
}

type state struct {
	minutes  int
	position vec2
}

type valley struct {
	start  vec2
	goal   vec2
	width  int
	height int
	up     map[vec2]bool
	down   map[vec2]bool
	left   map[vec2]bool
	right  map[vec2]bool
}

func parse() valley {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)
	}

	var start vec2
	for x, char := range lines[0] {
		if char == '.' {
			start = vec2{x - 1, -1}
		}
	}

	var goal vec2
	for x, char := range lines[len(lines)-1] {
		if char == '.' {
			goal = vec2{x - 1, len(lines) - 2}
		}
	}

	valley := valley{
		start,
		goal,
		len(lines[0]) - 2,
		len(lines) - 2,
		make(map[vec2]bool),
		make(map[vec2]bool),
		make(map[vec2]bool),
		make(map[vec2]bool),
	}
	for y, line := range lines {
		for x, char := range line {
			switch char {
			case '^':
				valley.up[vec2{x - 1, y - 1}] = true
			case 'v':
				valley.down[vec2{x - 1, y - 1}] = true
			case '<':
				valley.left[vec2{x - 1, y - 1}] = true
			case '>':
				valley.right[vec2{x - 1, y - 1}] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return valley
}

func safePosition(position vec2, valley valley, minutes int) bool {
	if (position.x < 0 || position.x >= valley.width || position.y < 0 || position.y >= valley.height) && position != valley.start && position != valley.goal {
		return false
	}

	invertedUp := position.y + minutes
	invertedDown := position.y - minutes
	for invertedDown < 0 {
		invertedDown += valley.height
	}
	invertedLeft := position.x + minutes
	invertedRight := position.x - minutes
	for invertedRight < 0 {
		invertedRight += valley.width
	}

	if valley.up[vec2{position.x, invertedUp % valley.height}] {
		return false
	}

	if valley.down[vec2{position.x, invertedDown % valley.height}] {
		return false
	}

	if valley.left[vec2{invertedLeft % valley.width, position.y}] {
		return false
	}

	if valley.right[vec2{invertedRight % valley.width, position.y}] {
		return false
	}

	return true
}

func fastestPath(valley valley, minutes int) int {
	start := state{minutes, valley.start}
	visited := make(map[state]bool)
	visited[start] = true
	queue := []state{start}

	rounds := 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		rounds++

		if current.position == valley.goal {
			return current.minutes
		}

		nextPosition := vec2{current.position.x, current.position.y - 1} // up
		if safePosition(nextPosition, valley, current.minutes+1) {
			nextState := state{current.minutes + 1, nextPosition}
			pseudoState := state{nextState.minutes % (valley.width * valley.height), nextState.position}
			if !visited[pseudoState] {
				visited[pseudoState] = true
				queue = append(queue, nextState)
			}
		}

		nextPosition = vec2{current.position.x, current.position.y + 1} // down
		if safePosition(nextPosition, valley, current.minutes+1) {
			nextState := state{current.minutes + 1, nextPosition}
			pseudoState := state{nextState.minutes % (valley.width * valley.height), nextState.position}
			if !visited[pseudoState] {
				visited[pseudoState] = true
				queue = append(queue, state{current.minutes + 1, nextPosition})
			}
		}

		nextPosition = vec2{current.position.x - 1, current.position.y} // left
		if safePosition(nextPosition, valley, current.minutes+1) {
			nextState := state{current.minutes + 1, nextPosition}
			pseudoState := state{nextState.minutes % (valley.width * valley.height), nextState.position}
			if !visited[pseudoState] {
				visited[pseudoState] = true
				queue = append(queue, state{current.minutes + 1, nextPosition})
			}
		}

		nextPosition = vec2{current.position.x + 1, current.position.y} // right
		if safePosition(nextPosition, valley, current.minutes+1) {
			nextState := state{current.minutes + 1, nextPosition}
			pseudoState := state{nextState.minutes % (valley.width * valley.height), nextState.position}
			if !visited[pseudoState] {
				visited[pseudoState] = true
				queue = append(queue, state{current.minutes + 1, nextPosition})
			}
		}

		nextPosition = current.position // wait
		if safePosition(nextPosition, valley, current.minutes+1) {
			nextState := state{current.minutes + 1, nextPosition}
			pseudoState := state{nextState.minutes % (valley.width * valley.height), nextState.position}
			if !visited[pseudoState] {
				visited[pseudoState] = true
				queue = append(queue, state{current.minutes + 1, nextPosition})
			}
		}
	}

	return -1
}

func part1(valley valley) int {
	return fastestPath(valley, 0)
}

func part2(valley valley) int {
	trip1 := fastestPath(valley, 0)
	valley.start, valley.goal = valley.goal, valley.start
	trip2 := fastestPath(valley, trip1)
	valley.start, valley.goal = valley.goal, valley.start
	trip3 := fastestPath(valley, trip2)
	return trip3
}

func main() {
	valley := parse()
	fmt.Println(part1(valley))
	fmt.Println(part2(valley))
}
