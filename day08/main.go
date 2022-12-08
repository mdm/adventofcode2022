package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parse() [][]int {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var forest [][]int

	for scanner.Scan() {
		line := scanner.Text()

		var trees []int
		for _, tree := range strings.Split(line, "") {
			height, err := strconv.Atoi(tree)
			if err != nil {
				log.Fatal(err)
			}

			trees = append(trees, height)
		}

		forest = append(forest, trees)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return forest
}

func part1(forest [][]int) int {
	h := len(forest)
	w := len(forest[0])

	count := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			visible := true
			for i := x; i > 0; i-- {
				if forest[y][x] <= forest[y][i-1] {
					visible = false
					break
				}
			}
			if visible {
				count++
				continue
			}

			visible = true
			for i := x; i < w-1; i++ {
				if forest[y][x] <= forest[y][i+1] {
					visible = false
					break
				}
			}
			if visible {
				count++
				continue
			}

			visible = true
			for i := y; i > 0; i-- {
				if forest[y][x] <= forest[i-1][x] {
					visible = false
					break
				}
			}
			if visible {
				count++
				continue
			}

			visible = true
			for i := y; i < h-1; i++ {
				if forest[y][x] <= forest[i+1][x] {
					visible = false
					break
				}
			}
			if visible {
				count++
				continue
			}
		}
	}

	return count
}

func part2(forest [][]int) int {
	h := len(forest)
	w := len(forest[0])

	maxScore := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			left := 0
			for i := x; i > 0; i-- {
				left++
				if forest[y][x] <= forest[y][i-1] {
					break
				}
			}

			right := 0
			for i := x; i < w-1; i++ {
				right++
				if forest[y][x] <= forest[y][i+1] {
					break
				}
			}

			up := 0
			for i := y; i > 0; i-- {
				up++
				if forest[y][x] <= forest[i-1][x] {
					break
				}
			}

			down := 0
			for i := y; i < h-1; i++ {
				down++
				if forest[y][x] <= forest[i+1][x] {
					break
				}
			}

			score := left * right * up * down
			if score > maxScore {
				maxScore = score
			}
		}
	}

	return maxScore
}

func main() {
	forest := parse()
	fmt.Println(part1(forest))
	fmt.Println(part2(forest))
}
