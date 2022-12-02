package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func loads() (loads []int) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	load := 0
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			loads = append(loads, load)
			load = 0
			continue
		}

		item, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}

		load += item
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return loads
}

func part1(loads []int) int {
	max := 0
	for _, load := range loads {
		if load > max {
			max = load
		}
	}

	return max
}

func part2(loads []int) int {
	sort.Sort(sort.Reverse(sort.IntSlice(loads)))

	top3 := 0
	for _, load := range loads[:3] {
		top3 += load
	}

	return top3
}

func main() {
	loads := loads()
	fmt.Println(part1(loads))
	fmt.Println(part2(loads))
}
