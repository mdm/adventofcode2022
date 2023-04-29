package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type node struct {
	value    int64
	previous *node
	next     *node
}

func parse() (nodes []*node, zero *node) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var previous *node
	for scanner.Scan() {
		line := scanner.Text()
		value, _ := strconv.Atoi(line)

		node := &node{int64(value), previous, nil}

		if previous != nil {
			previous.next = node
		}

		if node.value == 0 {
			zero = node
		}

		nodes = append(nodes, node)
		previous = node
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	nodes[0].previous = nodes[len(nodes)-1]
	nodes[len(nodes)-1].next = nodes[0]

	return nodes, zero
}

func print(nodes []*node) {
	cursor := nodes[0]
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("%d ", cursor.value)
		cursor = cursor.next
	}
	fmt.Println()
}

func mix(nodes []*node) {
	for _, n := range nodes {
		if n.value >= 0 {
			for i := int64(0); i < (n.value % (int64(len(nodes)) - int64(1))); i++ {
				oldPrevious := n.previous
				oldNext := n.next
				oldNextNext := n.next.next

				// unlink n
				oldPrevious.next = oldNext
				oldNext.previous = oldPrevious

				// relink n
				oldNext.next = n
				n.previous = oldNext
				oldNextNext.previous = n
				n.next = oldNextNext
			}
		} else {
			for i := int64(0); i < -(n.value % (int64(len(nodes)) - int64(1))); i++ {
				oldPrevious := n.previous
				oldPreviousPrevious := n.previous.previous
				oldNext := n.next

				// unlink n
				oldPrevious.next = oldNext
				oldNext.previous = oldPrevious

				oldPrevious.previous = n
				n.next = oldPrevious
				oldPreviousPrevious.next = n
				n.previous = oldPreviousPrevious
			}
		}
	}
}

func part1(nodes []*node, zero *node) int64 {
	mix(nodes)

	sum := int64(0)
	cursor := zero
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			cursor = cursor.next
		}

		sum += cursor.value
	}

	return sum
}

func part2(nodes []*node, zero *node) int64 {
	for _, n := range nodes {
		n.value *= 811_589_153
	}

	for i := 0; i < 10; i++ {
		mix(nodes)
	}

	sum := int64(0)
	cursor := zero
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			cursor = cursor.next
		}

		sum += cursor.value
	}

	return sum
}

func main() {
	nodes, zero := parse()
	fmt.Println(part1(nodes, zero))
	nodes, zero = parse()
	fmt.Println(part2(nodes, zero))
}
