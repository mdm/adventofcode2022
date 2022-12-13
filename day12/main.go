package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type position struct {
	x int
	y int
}

func parse() (locations [][]int, start position, destination position) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	y := 0
	for scanner.Scan() {
		line := scanner.Text()

		var elevations []int
		for x, char := range strings.Split(line, "") {
			var elevation int
			switch char {
			case "S":
				start = position{x, y}
				elevation = 0
			case "E":
				destination = position{x, y}
				elevation = 25
			default:
				elevation = int(char[0]) - int('a')
			}

			elevations = append(elevations, elevation)
		}
		locations = append(locations, elevations)
		y++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return locations, start, destination
}

func part1(locations [][]int, start position, destination position) int {
	visited := make(map[position]bool)
	steps := make(map[position]int)
	for y, row := range locations {
		for x := range row {
			pos := position{x, y}
			visited[pos] = false
			steps[pos] = 0
		}
	}

	var queue []position = []position{start}
	visited[start] = true
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == destination {
			return steps[destination]
		}

		// up
		next := position{current.x, current.y - 1}
		if next.y >= 0 && !visited[next] && locations[next.y][next.x] <= locations[current.y][current.x]+1 {
			queue = append(queue, next)
			visited[next] = true
			steps[next] = steps[current] + 1
		}

		// right
		next = position{current.x + 1, current.y}
		if next.x <= len(locations[0])-1 && !visited[next] && locations[next.y][next.x] <= locations[current.y][current.x]+1 {
			queue = append(queue, next)
			visited[next] = true
			steps[next] = steps[current] + 1
		}

		// down
		next = position{current.x, current.y + 1}
		if next.y <= len(locations)-1 && !visited[next] && locations[next.y][next.x] <= locations[current.y][current.x]+1 {
			queue = append(queue, next)
			visited[next] = true
			steps[next] = steps[current] + 1
		}

		// left
		next = position{current.x - 1, current.y}
		if next.x >= 0 && !visited[next] && locations[next.y][next.x] <= locations[current.y][current.x]+1 {
			queue = append(queue, next)
			visited[next] = true
			steps[next] = steps[current] + 1
		}
	}

	return -1
}

func part2(locations [][]int, destination position) int {
	min := len(locations) * len(locations[0])

	for y, row := range locations {
		for x, elevation := range row {
			if elevation == 0 {
				steps := part1(locations, position{x, y}, destination)

				if steps > -1 && steps < min {
					min = steps
				}
			}
		}
	}

	return min
}

func main() {
	locations, start, destination := parse()
	fmt.Println(part1(locations, start, destination))
	fmt.Println(part2(locations, destination))
}
