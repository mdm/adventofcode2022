package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type rucksack struct {
	first  []rune
	second []rune
	all    []rune
}

func rucksacks() (rucksacks []rucksack) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		size := len(line) / 2

		first := []rune(line[:size])
		sort.Slice(first, func(i int, j int) bool { return first[i] < first[j] })

		second := []rune(line[size:])
		sort.Slice(second, func(i int, j int) bool { return second[i] < second[j] })

		all := []rune(line)
		sort.Slice(all, func(i int, j int) bool { return all[i] < all[j] })

		rucksack := rucksack{first, second, all}
		rucksacks = append(rucksacks, rucksack)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return rucksacks
}

func part1(rucksacks []rucksack) int {
	sum := 0
	for _, rucksack := range rucksacks {
		var priority int
		f := 0
		s := 0
		for f < len(rucksack.first) && s < len(rucksack.second) {
			if rucksack.first[f] == rucksack.second[s] {
				if int(rucksack.first[f]) >= int('a') && int(rucksack.first[f]) <= int('z') {
					priority = int(rucksack.first[f]) - int('a') + 1
				}
				if int(rucksack.first[f]) >= int('A') && int(rucksack.first[f]) <= int('Z') {
					priority = int(rucksack.first[f]) - int('A') + 27
				}

				break
			} else if rucksack.first[f] < rucksack.second[s] {
				f++
			} else if rucksack.first[f] > rucksack.second[s] {
				s++
			}
		}

		sum += priority
	}

	return sum
}

func part2(rucksacks []rucksack) int {
	sum := 0

	for start := 0; start < len(rucksacks); start += 3 {
		var priority int
		f := 0
		s := 0
		t := 0
		for f < len(rucksacks[start].all) && s < len(rucksacks[start+1].all) && t < len(rucksacks[start+2].all) {
			items := []int{int(rucksacks[start].all[f]), int(rucksacks[start+1].all[s]), int(rucksacks[start+2].all[t])}
			sort.Sort(sort.IntSlice(items))
			min := items[0]

			if int(rucksacks[start].all[f]) == min && int(rucksacks[start+1].all[s]) == min && int(rucksacks[start+2].all[t]) == min {
				if int(rucksacks[start].all[f]) >= int('a') && int(rucksacks[start].all[f]) <= int('z') {
					priority = int(rucksacks[start].all[f]) - int('a') + 1
				}
				if int(rucksacks[start].all[f]) >= int('A') && int(rucksacks[start].all[f]) <= int('Z') {
					priority = int(rucksacks[start].all[f]) - int('A') + 27
				}

				break
			} else {
				if int(rucksacks[start].all[f]) == min {
					f++
				}
				if int(rucksacks[start+1].all[s]) == min {
					s++
				}
				if int(rucksacks[start+2].all[t]) == min {
					t++
				}
			}
		}

		sum += priority
	}

	return sum
}

func main() {
	rucksacks := rucksacks()
	fmt.Println(part1(rucksacks))
	fmt.Println(part2(rucksacks))
}
