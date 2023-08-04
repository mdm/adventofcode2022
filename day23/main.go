package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"

	"golang.org/x/exp/slices"
)

type vec2 struct {
	x int
	y int
}

func parse() (elves map[vec2]int) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	elves = make(map[vec2]int)
	id := 0
	y := 0
	for scanner.Scan() {
		line := scanner.Text()

		for x, char := range line {
			if char == '#' {
				elves[vec2{x, y}] = id
				id += 1
			}
		}

		y += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return elves
}

func step(prevElves map[vec2]int, modifier int) (map[vec2]int, int) {
	nextElves := make(map[vec2]int)
	moves := make(map[vec2][]int)
	stayed := 0
	for prevPosition, elf := range prevElves {
		stay := true
		for _, offset := range []vec2{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}} {
			neighbor := vec2{prevPosition.x + offset.x, prevPosition.y + offset.y}
			_, ok := prevElves[neighbor]

			if ok {
				stay = false
				break
			}
		}

		if stay {
			moves[prevPosition] = append(moves[prevPosition], elf)
			stayed++
			continue
		}

		move, nextPosition := checkDirection(prevElves, prevPosition, (0+modifier)%4)
		if move {
			moves[nextPosition] = append(moves[nextPosition], elf)
			continue
		}

		move, nextPosition = checkDirection(prevElves, prevPosition, (1+modifier)%4)
		if move {
			moves[nextPosition] = append(moves[nextPosition], elf)
			continue
		}

		move, nextPosition = checkDirection(prevElves, prevPosition, (2+modifier)%4)
		if move {
			moves[nextPosition] = append(moves[nextPosition], elf)
			continue
		}

		move, nextPosition = checkDirection(prevElves, prevPosition, (3+modifier)%4)
		if move {
			moves[nextPosition] = append(moves[nextPosition], elf)
			continue
		}

		moves[prevPosition] = append(moves[prevPosition], elf)
		stayed++
	}

	moved := 0
	for nextPosition, elves := range moves {
		if len(elves) == 1 {
			nextElves[nextPosition] = elves[0]
			moved++
		} else {
			for prevPosition, elf := range prevElves {
				if slices.Contains(elves, elf) {
					nextElves[prevPosition] = elf
					stayed++
				}
			}
		}
	}

	return nextElves, moved - stayed
}

func checkDirection(prevElves map[vec2]int, prevPosition vec2, direction int) (bool, vec2) {
	switch direction {
	case 0:
		north := true
		for _, offset := range []vec2{{-1, -1}, {0, -1}, {1, -1}} {
			neighbor := vec2{prevPosition.x + offset.x, prevPosition.y + offset.y}
			_, ok := prevElves[neighbor]

			if ok {
				north = false
				break
			}
		}

		if north {
			nextPosition := vec2{prevPosition.x, prevPosition.y - 1}
			return true, nextPosition
		}

	case 1:
		south := true
		for _, offset := range []vec2{{-1, 1}, {0, 1}, {1, 1}} {
			neighbor := vec2{prevPosition.x + offset.x, prevPosition.y + offset.y}
			_, ok := prevElves[neighbor]

			if ok {
				south = false
				break
			}
		}

		if south {
			nextPosition := vec2{prevPosition.x, prevPosition.y + 1}
			return true, nextPosition
		}

	case 2:
		west := true
		for _, offset := range []vec2{{-1, -1}, {-1, 0}, {-1, 1}} {
			neighbor := vec2{prevPosition.x + offset.x, prevPosition.y + offset.y}
			_, ok := prevElves[neighbor]

			if ok {
				west = false
				break
			}
		}

		if west {
			nextPosition := vec2{prevPosition.x - 1, prevPosition.y}
			return true, nextPosition
		}

	case 3:
		east := true
		for _, offset := range []vec2{{1, -1}, {1, 0}, {1, 1}} {
			neighbor := vec2{prevPosition.x + offset.x, prevPosition.y + offset.y}
			_, ok := prevElves[neighbor]

			if ok {
				east = false
				break
			}
		}

		if east {
			nextPosition := vec2{prevPosition.x + 1, prevPosition.y}
			return true, nextPosition
		}
	}

	return false, vec2{}
}

func part1(elves map[vec2]int) int {
	for i := 0; i < 10; i++ {
		nextElves, _ := step(elves, i)
		elves = nextElves
	}

	minX := math.MaxInt32
	maxX := math.MinInt32
	minY := math.MaxInt32
	maxY := math.MinInt32
	for position := range elves {
		if position.x < minX {
			minX = position.x
		}

		if position.x > maxX {
			maxX = position.x
		}

		if position.y < minY {
			minY = position.y
		}

		if position.y > maxY {
			maxY = position.y
		}
	}

	width := maxX - minX + 1
	height := maxY - minY + 1

	return width*height - len(elves)
}

func part2(elves map[vec2]int) int {
	round := 0
	for {
		nextElves, moved := step(elves, round)
		if moved == 0 {
			break
		}
		elves = nextElves
		round++
	}

	return round + 1
}

func main() {
	elves := parse()
	fmt.Println(part1(elves))
	fmt.Println(part2(elves))
}
