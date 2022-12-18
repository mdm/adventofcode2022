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

func stateHash(location string, minutes int, open map[string]int) string {
	var ids []string
	for id, _ := range open {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	hash := location + fmt.Sprint(minutes)
	for _, id := range ids {
		hash += id + fmt.Sprint(open[id])
	}

	return hash
}

func findMaxFlow(valves map[string]valve, location string, open map[string]int, minutes int, memo map[string]int) int {
	m, ok := memo[stateHash(location, minutes, open)]
	if ok {
		return m
	}

	if minutes == 30 {
		memo[stateHash(location, minutes, open)] = 0
		return 0
	}

	max := 0

	_, opened := open[location]
	if !opened && valves[location].flow > 0 {
		open[location] = minutes
		max = findMaxFlow(valves, location, open, minutes+1, memo)
		delete(open, location)
	}

	for _, next := range valves[location].tunnels {
		flow := findMaxFlow(valves, next, open, minutes+1, memo)

		if flow > max {
			max = flow
		}
	}

	for id, _ := range open {
		max += valves[id].flow
	}

	memo[stateHash(location, minutes, open)] = max
	return max
}

func part1(valves map[string]valve) int {
	open := make(map[string]int)
	location := "AA"
	memo := make(map[string]int)

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
