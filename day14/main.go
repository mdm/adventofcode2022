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

type rock struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func parse() (grid [][]rune, xOffset int, yOffset int) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var rocks []rock
	xMin := math.MaxInt
	xMax := 0
	yMin := math.MaxInt
	yMax := 0
	for scanner.Scan() {
		line := scanner.Text()
		coords := strings.Split(line, " -> ")

		var x1, y1, x2, y2 int
		i := 0
		for _, coord := range coords {
			parts := strings.Split(coord, ",")

			x, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Fatal(err)
			}

			if i != 0 {
				x2 = x
				y2 = y

				rocks = append(rocks, rock{x1, y1, x2, y2})
			}

			x1 = x
			y1 = y

			if x < xMin {
				xMin = x
			}

			if x > xMax {
				xMax = x
			}

			if y < yMin {
				yMin = y
			}

			if y > yMax {
				yMax = y
			}

			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if 500-(yMax+1) < xMin {
		xMin = 500 - (yMax + 1)
	}

	if 500+(yMax+1) > xMax {
		xMax = 500 + (yMax + 1)
	}

	return makeGrid(rocks, xMin-1, xMax+1, 0, yMax+1), xMin - 1, 0
}

func makeGrid(rocks []rock, xMin int, xMax int, yMin int, yMax int) (grid [][]rune) {
	width := xMax - xMin + 1
	height := yMax - yMin + 1

	for y := 0; y < height; y++ {
		var row []rune
		for x := 0; x < width; x++ {
			row = append(row, '.')
		}
		grid = append(grid, row)
	}

	for _, rock := range rocks {
		if rock.x1 == rock.x2 {
			// vertical
			if rock.y2 < rock.y1 {
				rock.y2, rock.y1 = rock.y1, rock.y2
			}
			for y := rock.y1; y <= rock.y2; y++ {
				grid[y-yMin][rock.x1-xMin] = '#'
			}
		} else {
			// horizontal
			if rock.x2 < rock.x1 {
				rock.x2, rock.x1 = rock.x1, rock.x2
			}
			for x := rock.x1; x <= rock.x2; x++ {
				grid[rock.y1-yMin][x-xMin] = '#'
			}
		}
	}

	return grid
}

func printGrid(grid [][]rune) {
	for _, row := range grid {
		fmt.Println(string(row))
	}
}

func part1(grid [][]rune, xOffset int, yOffset int) int {
	count := 0
	for {
		sandX := 500 - xOffset
		sandY := 0 - yOffset

		for {
			if sandY+1 == len(grid) {
				// printGrid(grid)
				return count
			}

			if grid[sandY+1][sandX] == '.' {
				sandY += 1
				continue
			}

			if grid[sandY+1][sandX-1] == '.' {
				sandX -= 1
				sandY += 1
				continue
			}

			if grid[sandY+1][sandX+1] == '.' {
				sandX += 1
				sandY += 1
				continue
			}

			grid[sandY][sandX] = 'o'
			count++
			break
		}
	}
}

func part2(grid [][]rune, xOffset int, yOffset int) int {
	count := 0
	for {
		sandX := 500 - xOffset
		sandY := 0 - yOffset

		for {
			if sandY+1 == len(grid) {
				grid[sandY][sandX] = 'o'
				count++
				break
			}

			if grid[sandY+1][sandX] == '.' {
				sandY += 1
				continue
			}

			if grid[sandY+1][sandX-1] == '.' {
				sandX -= 1
				sandY += 1
				continue
			}

			if grid[sandY+1][sandX+1] == '.' {
				sandX += 1
				sandY += 1
				continue
			}

			grid[sandY][sandX] = 'o'
			count++
			break
		}

		if sandX == 500-xOffset && sandY == 0-yOffset {
			// printGrid(grid)
			return count
		}
	}
}

func main() {
	grid, xOffset, yOffset := parse()
	answer1 := part1(grid, xOffset, yOffset)
	fmt.Println(answer1)
	fmt.Println(answer1 + part2(grid, xOffset, yOffset))
}
