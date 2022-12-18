package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type valve struct {
	id      string
	flow    int
	tunnels []string
}

func parse() (valves map[string]valve) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	regex := regexp.MustCompile(`Valve ([A-Z][A-Z]) has flow rate=(\d+); tunnels? leads? to valves? (.+)`)

	scanner := bufio.NewScanner(file)

	valves = make(map[string]valve)

	for scanner.Scan() {
		line := scanner.Text()
		match := regex.FindStringSubmatch(line)

		id := match[1]

		flow, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}

		var tunnels = strings.Split(match[3], ", ")

		valves[id] = valve{id, flow, tunnels}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	dotfile, err := os.Create(os.Args[1] + ".dot")
	if err != nil {
		log.Fatal(err)
	}
	defer dotfile.Close()

	dotfile.WriteString("digraph cave {\n")
	for id, valve := range valves {
		for _, tunnel := range valve.tunnels {
			dotfile.WriteString(fmt.Sprint(id) + " -> " + fmt.Sprint(tunnel) + "\n")
			if valve.flow > 0 {
				dotfile.WriteString(fmt.Sprint(id) + " [shape=box];")
			}
		}
	}
	dotfile.WriteString("}\n")

	return valves
}

func stateHash(location string, open map[string]bool) string {
	var ids []string
	for id, check := range open {
		if check {
			ids = append(ids, id)
		}
	}

	sort.Strings(ids)

	hash := location
	for _, id := range ids {
		hash += id
	}

	return hash
}

func findMaxFlow(valves map[string]valve, location string, open map[string]bool, minutes int, memo map[string]bool) int {
	if minutes == 30 {
		// fmt.Println(open)
		return 0
	}

	max := 0

	if !open[location] && valves[location].flow > 0 {
		// fmt.Println("open", location)
		open[location] = true
		// fmt.Println(location, minutes)
		max = findMaxFlow(valves, location, open, minutes+1, memo)
		open[location] = false
	} else {
		for id, check := range open {
			if check {
				max += (30 - minutes - 1) * valves[id].flow
			}
		}
	}

	// fmt.Println("skip", location)
	for _, next := range valves[location].tunnels {
		flow := 0
		if memo[stateHash(next, open)] {
			continue
		}

		memo[stateHash(next, open)] = true
		flow = findMaxFlow(valves, next, open, minutes+1, memo)

		if flow > max {
			max = flow
		}
	}

	for id, check := range open {
		if check {
			max += valves[id].flow
		}
	}

	return max
}

func part1(valves map[string]valve) int {
	fmt.Println(valves)

	open := make(map[string]bool)
	location := "AA"
	memo := make(map[string]bool)

	return findMaxFlow(valves, location, open, 0, memo)
}

func part2(valves map[string]valve) int {
	return 0
}

func main() {
	valves := parse()
	fmt.Println(part1(valves))
	fmt.Println(part2(valves))
}
