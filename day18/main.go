package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type cube struct {
	x int
	y int
	z int
}

func parse() (droplet map[cube]bool) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	droplet = make(map[cube]bool)

	for scanner.Scan() {
		line := scanner.Text()
		coords := strings.Split(line, ",")

		x, err := strconv.Atoi(coords[0])
		if err != nil {
			log.Fatal(err)
		}

		y, err := strconv.Atoi(coords[1])
		if err != nil {
			log.Fatal(err)
		}

		z, err := strconv.Atoi(coords[2])
		if err != nil {
			log.Fatal(err)
		}

		droplet[cube{x, y, z}] = true
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return droplet
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func part1(droplet map[cube]bool) int {
	surface := 0

	neighbors := []cube{{1, 0, 0}, {-1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1}}

	for c, _ := range droplet {
		for _, n := range neighbors {
			if !droplet[cube{c.x + n.x, c.y + n.y, c.z + n.z}] {
				surface++
			}
		}
	}

	return surface
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func part2(droplet map[cube]bool) int {
	surface := 0

	neighbors := []cube{{1, 0, 0}, {-1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1}}

	min := cube{math.MaxInt, math.MaxInt, math.MaxInt}
	max := cube{0, 0, 0}
	for c, _ := range droplet {
		if c.x < min.x {
			min.x = c.x
		}
		if c.y < min.y {
			min.y = c.y
		}
		if c.z < min.z {
			min.z = c.z
		}

		if c.x > max.x {
			max.x = c.x
		}
		if c.y > max.y {
			max.y = c.y
		}
		if c.z > max.z {
			max.z = c.z
		}
	}

	min.x--
	min.y--
	min.z--

	max.x++
	max.y++
	max.z++

	queue := []cube{min}
	visited := make(map[cube]bool)
	visited[min] = true

	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]

		for _, n := range neighbors {
			neighbor := cube{c.x + n.x, c.y + n.y, c.z + n.z}
			if droplet[neighbor] {
				surface++
			} else {
				if !visited[neighbor] && neighbor.x >= min.x && neighbor.x <= max.x && neighbor.y >= min.y && neighbor.y <= max.y && neighbor.z >= min.z && neighbor.z <= max.z {
					queue = append(queue, neighbor)
					visited[neighbor] = true
				}
			}
		}
	}

	return surface
}

func main() {
	droplet := parse()
	fmt.Println(part1(droplet))
	fmt.Println(part2(droplet))
}
