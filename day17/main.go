package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func printMountain(mountain []position) {
	height := 0
	for _, m := range mountain {
		if m.y > height {
			height = m.y
		}
	}

	fmt.Println(mountain)

	for y := height; y >= 0; y-- {
		row := []rune("|.......|")
		for _, m := range mountain {
			if m.y == y {
				row[m.x+1] = '#'
			}
		}
		fmt.Println(string(row))
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
		// fmt.Println(i)
		pos := position{2, height + 4}
		step := 0
		rock := rocks[r]
		r++
		r = r % len(rocks)

		for {
			oldPos := pos
			switch step % 2 {
			case 0:
				// fmt.Println("JET", pos)
				pos.x += jet[j]
				j++
				j = j % len(jet)
			case 1:
				// fmt.Println("DROP", pos)
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

					// printMountain(mountain)
					// fmt.Println()

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
	mountain := make(map[position]bool)
	for i := 0; i < 7; i++ {
		mountain[position{i, 0}] = true
	}
	r := 0
	j := 0

	// count := 1_000_000_000_000
	period := 5 * len(jet)

	var results, diffs []int
	height, newMountain, newR, newJ := dropRocks(mountain, r, j, jet, period)
	mountain = newMountain
	r = newR
	j = newJ
	results = append(results, height)
	fmt.Println("BAAM", results[0])
	for i := 0; i < 10_000; i++ {
		height, newMountain, newR, newJ := dropRocks(mountain, r, j, jet, period)
		mountain = newMountain
		r = newR
		j = newJ
		results = append(results, height)
		diff := results[i+1] - results[i]
		diffs = append(diffs, diff)
		fmt.Println(diff)

		cycle := 0
		for clen := 1; clen <= (len(diffs))/3; clen++ {
			valid := true
			for c := 0; c < clen; c++ {
				if diffs[len(diffs)-1-c] != diffs[len(diffs)-1-c-clen] || diffs[len(diffs)-1-c] != diffs[len(diffs)-1-c-2*clen] {
					valid = false
					break
				}
			}

			if valid {
				cycle = clen
				break
			}
		}

		if cycle > 0 {
			fmt.Println(i, cycle)
			break
		}
	}

	// 300
	// 306
	// 303
	// 303
	// 301
	// 306
	// 301

	return 0
}

func main() {
	jet := parse()
	fmt.Println(part1(jet))
	fmt.Println("------")
	fmt.Println(part2(jet))
}
