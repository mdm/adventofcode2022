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

type resource struct {
	kind   string
	amount int
}

type recipe struct {
	output string
	inputs []resource
}

func parse() (blueprints [][]recipe) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	regex := regexp.MustCompile(` Each (\w+) robot costs (\d+) (\w+)(:? and (\d+) (\w+))*`)

	for scanner.Scan() {
		line := scanner.Text()
		blueprint := strings.Split(line, ":")[1]

		var recipes []recipe
		for _, r := range strings.Split(blueprint, ".") {
			if len(r) == 0 {
				continue
			}

			match := regex.FindStringSubmatch(r)

			output := match[1]

			var inputs []resource
			for i := 2; i < len(match); i += 3 {
				if len(match[i]) == 0 {
					continue
				}

				amount, err := strconv.Atoi(match[i])
				if err != nil {
					log.Fatal(err)
				}

				input := resource{match[i+1], amount}
				inputs = append(inputs, input)
			}

			recipes = append(recipes, recipe{output, inputs})
		}

		blueprints = append(blueprints, recipes)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return blueprints
}

func stateHash(minutes int, resources map[string]int, robots map[string]int) string {
	hash := fmt.Sprint(minutes)

	var kinds []string
	for kind := range robots {
		kinds = append(kinds, kind)
	}

	sort.Strings(kinds)

	for _, kind := range kinds {
		hash += kind + fmt.Sprint(robots[kind]) + "/" + fmt.Sprint(resources[kind])
	}

	return hash
}

func optimizeBlueprint(minutes int, limit int, resources map[string]int, robots map[string]int, maxRobots map[string]int, blueprint []recipe, memo map[string]int) int {
	hash := stateHash(minutes, resources, robots)
	m, ok := memo[hash]
	if ok {
		return m
	}

	if minutes == limit {
		return resources["geode"]
	}

	newResources := make(map[string]int)
	for kind, amount := range resources {
		newResources[kind] = amount
	}

	for k, r := range robots {
		newResources[k] += r
	}

	for _, b := range blueprint {
		if b.output != "geode" {
			continue
		}
		build := true
		for _, i := range b.inputs {
			if resources[i.kind] < i.amount {
				build = false
				break
			}
		}

		if build {
			for _, i := range b.inputs {
				newResources[i.kind] -= i.amount
			}
			robots[b.output]++

			geodes := optimizeBlueprint(minutes+1, limit, newResources, robots, maxRobots, blueprint, memo)

			for _, i := range b.inputs {
				newResources[i.kind] += i.amount
			}
			robots[b.output]--

			memo[hash] = geodes
			return geodes
		}
	}

	max := 0
	for _, b := range blueprint {
		if b.output != "geode" && robots[b.output] >= maxRobots[b.output] {
			continue
		}

		build := true
		for _, i := range b.inputs {
			if resources[i.kind] < i.amount {
				build = false
				break
			}
		}

		if build {
			for _, i := range b.inputs {
				newResources[i.kind] -= i.amount
			}
			robots[b.output]++

			geodes := optimizeBlueprint(minutes+1, limit, newResources, robots, maxRobots, blueprint, memo)
			if geodes > max {
				max = geodes
			}

			for _, i := range b.inputs {
				newResources[i.kind] += i.amount
			}
			robots[b.output]--
		}
	}

	geodes := optimizeBlueprint(minutes+1, limit, newResources, robots, maxRobots, blueprint, memo)
	if geodes > max {
		max = geodes
	}

	memo[hash] = max
	return max
}

func part1(blueprints [][]recipe) int {
	sum := 0
	for i, blueprint := range blueprints {
		maxRobots := make(map[string]int)
		for _, b := range blueprint {
			for _, i := range b.inputs {
				if maxRobots[i.kind] < i.amount {
					maxRobots[i.kind] = i.amount
				}
			}
		}
		resources := make(map[string]int)
		robots := make(map[string]int)
		robots["ore"] = 1
		memo := make(map[string]int)
		geodes := optimizeBlueprint(0, 24, resources, robots, maxRobots, blueprint, memo)
		sum += (i + 1) * geodes
	}

	return sum
}

func part2(blueprints [][]recipe) int {
	product := 1
	for i, blueprint := range blueprints {
		if i == 3 {
			break
		}
		maxRobots := make(map[string]int)
		for _, b := range blueprint {
			for _, i := range b.inputs {
				if maxRobots[i.kind] < i.amount {
					maxRobots[i.kind] = i.amount
				}
			}
		}
		resources := make(map[string]int)
		robots := make(map[string]int)
		robots["ore"] = 1
		memo := make(map[string]int)
		geodes := optimizeBlueprint(0, 32, resources, robots, maxRobots, blueprint, memo)
		product *= geodes
	}

	return product
}

func main() {
	blueprints := parse()
	fmt.Println(part1(blueprints))
	fmt.Println(part2(blueprints))
}
