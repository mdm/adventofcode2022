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
			// fmt.Println(output)

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
			// fmt.Println(inputs)

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

func optimizeBlueprint(minutes int, resources map[string]int, robots map[string]int, blueprint []recipe, memo map[string]int) int {
	hash := stateHash(minutes, resources, robots)
	m, ok := memo[hash]
	if ok {
		return m
	}

	if minutes == 24 {
		return resources["geode"]
	}

	for k, r := range robots {
		resources[k] += r
	}

	max := optimizeBlueprint(minutes+1, resources, robots, blueprint, memo)

	for k, r := range robots {
		resources[k] -= r
	}

	for _, b := range blueprint {
		build := true
		for _, i := range b.inputs {
			if resources[i.kind] < i.amount {
				build = false
				break
			}
		}

		if build {
			for _, i := range b.inputs {
				resources[i.kind] -= i.amount
			}
			robots[b.output]++

			for k, r := range robots {
				resources[k] += r
			}

			geodes := optimizeBlueprint(minutes+1, resources, robots, blueprint, memo)
			if geodes > max {
				max = geodes
			}

			for k, r := range robots {
				resources[k] -= r
			}

			robots[b.output]--
			for _, i := range b.inputs {
				resources[i.kind] += i.amount
			}
		}
	}

	memo[hash] = max
	return max
}

func part1(blueprints [][]recipe) int {
	fmt.Println(blueprints)
	sum := 0

	for i, blueprint := range blueprints {
		resources := make(map[string]int)
		robots := make(map[string]int)
		robots["ore"] = 1
		memo := make(map[string]int)
		geodes := optimizeBlueprint(0, resources, robots, blueprint, memo)
		fmt.Println(i+1, geodes)
		sum += (i + 1) * geodes
	}

	return sum
}

func part2(blueprints [][]recipe) int {
	return 0
}

func main() {
	blueprints := parse()
	fmt.Println(part1(blueprints))
	fmt.Println(part2(blueprints))
}
