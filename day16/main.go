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
	id        string
	flow      int
	tunnels   []string
	distances map[string]int
}

type action struct {
	move     bool
	location string
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

		valves[id] = valve{id, flow, tunnels, nil}
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

func findDistances(valves map[string]valve, start string) map[string]int {
	visited := make(map[string]bool)
	steps := make(map[string]int)
	for v, _ := range valves {
		visited[v] = false
		steps[v] = 0
	}

	queue := []string{start}
	visited[start] = true
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, next := range valves[current].tunnels {
			if !visited[next] {
				queue = append(queue, next)
				visited[next] = true
				steps[next] = steps[current] + 1
			}
		}

	}

	return steps
}

func mapNetwork(valves map[string]valve) map[string]valve {
	mappedValves := make(map[string]valve)
	for id, v := range valves {
		mappedValves[id] = valve{id, v.flow, v.tunnels, findDistances(valves, id)}
	}

	return mappedValves
}

func findMaxFlow(valves map[string]valve, location string, open map[string]int, minutes int, limit int, memo map[string]int) int {
	max := 0
	progress := false
	for next, _ := range valves {
		_, opened := open[next]
		if !opened && valves[next].flow > 0 {
			eta := minutes + valves[location].distances[next] + 1
			if eta <= limit {
				progress = true
				open[next] = eta
				flow := findMaxFlow(valves, next, open, eta, limit, memo)
				delete(open, next)

				if flow > max {
					max = flow
				}
			}
		}
	}

	if !progress {
		flow := 0
		for id, minute := range open {
			flow += valves[id].flow * (limit - minute)
		}
		return flow
	}

	return max
}

func stateHash(locationMe string, locationElephant string, minutesMe int, minutesElephant int, open map[string]int) string {
	var ids []string
	for id, _ := range open {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	var hash string
	if locationMe < locationElephant {
		hash = locationMe + fmt.Sprint(minutesMe) + locationElephant + fmt.Sprint(minutesElephant)
	} else {
		hash = locationElephant + fmt.Sprint(minutesElephant) + locationMe + fmt.Sprint(minutesMe)
	}

	for _, id := range ids {
		hash += id + fmt.Sprint(open[id])
	}

	return hash
}

func findMaxFlow2(valves map[string]valve, ids []string, locationMe string, locationElephant string, open map[string]int, minutesMe int, minutesElephant int, limit int, memo map[string]int) int {
	hash := stateHash(locationMe, locationElephant, minutesMe, minutesElephant, open)
	m, ok := memo[hash]
	if ok {
		return m
	}

	max := 0
	bothProgess := false
	for i, nextMe := range ids {
		progressMe := false
		etaMe := 0
		_, openedMe := open[nextMe]
		if !openedMe && valves[nextMe].flow > 0 {
			etaMe = minutesMe + valves[locationMe].distances[nextMe] + 1
			if etaMe <= limit {
				progressMe = true
				open[nextMe] = etaMe
			}
		}

		if !progressMe {
			continue
		}

		for j, nextElephant := range ids {
			if i == j {
				continue
			}

			progressElephant := false
			etaElephant := 0
			_, openedElephant := open[nextElephant]
			if !openedElephant && valves[nextElephant].flow > 0 {
				etaElephant = minutesElephant + valves[locationElephant].distances[nextElephant] + 1
				if etaElephant <= limit {
					progressElephant = true
					open[nextElephant] = etaElephant
				}
			}

			if !progressElephant {
				continue
			}

			bothProgess = true
			flow := findMaxFlow2(valves, ids, nextMe, nextElephant, open, etaMe, etaElephant, limit, memo)
			if flow > max {
				max = flow
			}

			delete(open, nextElephant)
		}

		delete(open, nextMe)
	}

	if !bothProgess {
		flowMe := findMaxFlow(valves, locationMe, open, minutesMe, limit, nil)
		flowElephant := findMaxFlow(valves, locationElephant, open, minutesElephant, limit, nil)

		if flowMe > max {
			max = flowMe
		}
		if flowElephant > max {
			max = flowElephant
		}
	}

	memo[hash] = max
	return max
}

func part1(valves map[string]valve) int {
	valves = mapNetwork(valves)

	open := make(map[string]int)
	location := "AA"
	memo := make(map[string]int)

	flow := findMaxFlow(valves, location, open, 0, 30, memo)
	return flow
}

func part2(valves map[string]valve) int {
	valves = mapNetwork(valves)

	ids := make([]string, len(valves))
	i := 0
	for id, _ := range valves {
		ids[i] = id
		i++
	}

	open := make(map[string]int)
	locationMe := "AA"
	locationElephant := "AA"
	memo := make(map[string]int)

	flow := findMaxFlow2(valves, ids, locationMe, locationElephant, open, 0, 0, 26, memo)
	return flow
}

func main() {
	valves := parse()
	fmt.Println(part1(valves))
	fmt.Println(part2(valves))
}
