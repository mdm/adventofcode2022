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

func parse() (jet []int) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		line := scanner.Text()

		for _, char := range line {
			switch char {
			case '<':
				jet = append(jet, -1)
			case '>':
				jet = append(jet, 1)
			default:
				log.Fatal(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return jet
}

var rock0 = []position{{0, 0}, {1, 0}, {2, 0}, {3, 0}}
var rock1 = []position{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}}
var rock2 = []position{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}}
var rock3 = []position{{0, 0}, {0, 1}, {0, 2}, {0, 3}}
var rock4 = []position{{0, 0}, {1, 0}, {0, 1}, {1, 1}}

var rocks = [][]position{rock0, rock1, rock2, rock3, rock4}

func printMountain(mountain map[position]bool, height int) {
	for y := height; y >= 0; y-- {
		row := []rune("|.......|")
		for x := 0; x < 7; x++ {
			_, ok := mountain[position{x, y}]
			if ok {
				row[x+1] = '#'
			}

		}
		fmt.Print(string(row))
		fmt.Println(y)
	}
}

func dropRocks(mountain map[position]bool, r int, j int, jet []int, count int) (int, map[position]bool, int, int) {
	height := 0
	for m, _ := range mountain {
		if m.y > height {
			height = m.y
		}
	}

	for i := 0; i < count; i++ {
		pos := position{2, height + 4}
		step := 0
		rock := rocks[r]
		r++
		r = r % len(rocks)

		for {
			oldPos := pos
			switch step % 2 {
			case 0:
				pos.x += jet[j]
				j++
				j = j % len(jet)
			case 1:
				pos.y -= 1
			}

			crash := false
			for _, p := range rock {
				if p.x+pos.x < 0 || p.x+pos.x > 6 {
					crash = true
					break
				}
			}

			if !crash {
				for _, p := range rock {
					if mountain[position{p.x + pos.x, p.y + pos.y}] {
						crash = true
						break
					}
				}
			}

			if crash {
				pos = oldPos

				if step%2 == 1 {
					for _, p := range rock {
						mountain[position{p.x + pos.x, p.y + pos.y}] = true

						if p.y+pos.y > height {
							height = p.y + pos.y
						}
					}

					break
				}
			}

			step += 1
		}
	}

	return height, mountain, r, j
}

func part1(jet []int) int {
	mountain := make(map[position]bool)
	for i := 0; i < 7; i++ {
		mountain[position{i, 0}] = true
	}

	height, _, _, _ := dropRocks(mountain, 0, 0, jet, 2022)
	return height
}

func part2(jet []int) int {
	memo := make(map[string]int)
	rs := make(map[string]int)
	js := make(map[string]int)
	mountain := make(map[position]bool)
	var columns []int
	for i := 0; i < 7; i++ {
		mountain[position{i, 0}] = true
		columns = append(columns, 0)
	}

	hash := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(columns)), ","), "[]")
	memo[hash] = 0

	r := 0
	j := 0
	rs[hash] = r
	js[hash] = j

	dropped := 0
	height := 0
	heights := []int{0}
	for {
		columns = nil
		newHeight, newMountain, newR, newJ := dropRocks(mountain, r, j, jet, 1)
		height = newHeight
		mountain = newMountain
		r = newR
		j = newJ
		dropped++
		heights = append(heights, height)

		min := height
		for i := 0; i < 7; i++ {
			max := height
			for max >= 0 && !mountain[position{i, max}] {
				max--
			}

			if max < min {
				min = max
			}

			columns = append(columns, max)
		}

		for i := 0; i < 7; i++ {
			columns[i] -= min
		}

		hash = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(columns)), ","), "[]")

		_, ok := memo[hash]
		if ok && rs[hash] == r && js[hash] == j {
			break
		}

		memo[hash] = dropped
		rs[hash] = r
		js[hash] = j
	}

	quotient := (1_000_000_000_000 - memo[hash]) / (dropped - memo[hash])
	remainder := (1_000_000_000_000 - memo[hash]) % (dropped - memo[hash])
	return quotient*(heights[dropped]-heights[memo[hash]]) + heights[memo[hash]+remainder]
}

func main() {
	jet := parse()
	fmt.Println(part1(jet))
	fmt.Println(part2(jet))
}
