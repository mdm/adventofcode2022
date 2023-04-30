package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type node interface {
	evaluate(nodes map[string]node) int
	infer(nodes map[string]node, target int) int
	complete(nodes map[string]node) bool
}

type inner struct {
	self      string
	left      string
	right     string
	operation string
}

type leaf struct {
	self  string
	value int
}

func parse() (nodes map[string]node) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	nodes = make(map[string]node)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, ": ")
		name := parts[0]
		value, err := strconv.Atoi(parts[1])

		if err == nil {
			l := leaf{name, value}
			nodes[name] = l
		} else {
			parts = strings.Split(parts[1], " ")
			i := inner{name, parts[0], parts[2], parts[1]}
			nodes[name] = i
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nodes
}

func (i inner) evaluate(nodes map[string]node) int {
	switch i.operation {
	case "+":
		return nodes[i.left].evaluate(nodes) + nodes[i.right].evaluate(nodes)
	case "-":
		return nodes[i.left].evaluate(nodes) - nodes[i.right].evaluate(nodes)
	case "*":
		return nodes[i.left].evaluate(nodes) * nodes[i.right].evaluate(nodes)
	case "/":
		return nodes[i.left].evaluate(nodes) / nodes[i.right].evaluate(nodes)
	}

	return 0
}

func (i inner) infer(nodes map[string]node, target int) int {
	if nodes[i.left].complete(nodes) {
		left := nodes[i.left].evaluate(nodes)
		switch i.operation {
		case "+":
			// left + x = target
			return nodes[i.right].infer(nodes, target-left)
		case "-":
			// left - x = target
			return nodes[i.right].infer(nodes, left-target)
		case "*":
			// left * x = target
			return nodes[i.right].infer(nodes, target/left)
		case "/":
			// left / x = target
			return nodes[i.right].infer(nodes, left/target)
		}
	}

	if nodes[i.right].complete(nodes) {
		right := nodes[i.right].evaluate(nodes)
		switch i.operation {
		case "+":
			// x + right = target
			return nodes[i.left].infer(nodes, target-right)
		case "-":
			// x - right = target
			return nodes[i.left].infer(nodes, target+right)
		case "*":
			// x * right = target
			return nodes[i.left].infer(nodes, target/right)
		case "/":
			// x / right = target
			return nodes[i.left].infer(nodes, target*right)
		}
	}

	return 0
}

func (i inner) complete(nodes map[string]node) bool {
	return nodes[i.left].complete(nodes) && nodes[i.right].complete(nodes)
}

func (l leaf) evaluate(nodes map[string]node) int {
	return l.value
}

func (l leaf) infer(nodes map[string]node, target int) int {
	if l.self == "humn" {
		return target
	} else {
		return l.value
	}
}

func (l leaf) complete(nodes map[string]node) bool {
	return l.self != "humn"
}

func part1(nodes map[string]node) int {
	return nodes["root"].evaluate(nodes)
}

func part2(nodes map[string]node) int {
	left := nodes["root"].(inner).left
	right := nodes["root"].(inner).right

	if nodes[left].complete(nodes) {
		return nodes[right].infer(nodes, nodes[left].evaluate(nodes))
	}

	if nodes[right].complete(nodes) {
		return nodes[left].infer(nodes, nodes[right].evaluate(nodes))
	}

	return 0
}

func main() {
	nodes := parse()
	fmt.Println(part1(nodes))
	fmt.Println(part2(nodes))
}
