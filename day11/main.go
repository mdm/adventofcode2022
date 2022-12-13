package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type operand struct {
	old   bool
	value int
}

type monkey struct {
	items     []int
	operation string
	operands  []operand
	test      int
	success   int
	failure   int
}

func parse() []monkey {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var monkeys []monkey
	var m monkey

	phase := 0
	for scanner.Scan() {
		line := scanner.Text()

		switch phase {
		case 0:
		case 1:
			line = strings.Split(line, ": ")[1]
			items := strings.Split(line, ", ")

			for _, item := range items {
				worry, err := strconv.Atoi(item)
				if err != nil {
					log.Fatal(err)
				}

				m.items = append(m.items, worry)
			}
		case 2:
			line = strings.Split(line, " = ")[1]
			items := strings.Split(line, " ")

			m.operation = items[1]

			var op operand
			if items[0] == "old" {
				op.old = true
			} else {
				value, err := strconv.Atoi(items[0])
				if err != nil {
					log.Fatal(err)
				}

				op.value = value
			}
			m.operands = append(m.operands, op)

			op = operand{}
			if items[2] == "old" {
				op.old = true
			} else {
				value, err := strconv.Atoi(items[2])
				if err != nil {
					log.Fatal(err)
				}

				op.value = value
			}
			m.operands = append(m.operands, op)
		case 3:
			items := strings.Split(line, " ")
			test, err := strconv.Atoi(items[len(items)-1])
			if err != nil {
				log.Fatal(err)
			}
			m.test = test
		case 4:
			items := strings.Split(line, " ")
			success, err := strconv.Atoi(items[len(items)-1])
			if err != nil {
				log.Fatal(err)
			}
			m.success = success
		case 5:
			items := strings.Split(line, " ")
			failure, err := strconv.Atoi(items[len(items)-1])
			if err != nil {
				log.Fatal(err)
			}
			m.failure = failure
		case 6:
		}

		phase++

		if phase == 7 {
			monkeys = append(monkeys, m)
			m = monkey{}
			phase = 0
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	monkeys = append(monkeys, m)

	return monkeys
}

func part1(monkeys []monkey) int {
	var inspected []int
	for range monkeys {
		inspected = append(inspected, 0)
	}

	for round := 0; round < 20; round++ {
		for i, m := range monkeys {
			for _, item := range m.items {
				inspected[i]++

				var a int
				if m.operands[0].old {
					a = item
				} else {
					a = m.operands[0].value
				}

				var b int
				if m.operands[1].old {
					b = item
				} else {
					b = m.operands[1].value
				}

				switch m.operation {
				case "+":
					item = a + b
				case "*":
					item = a * b
				}

				item /= 3

				if item%m.test == 0 {
					monkeys[m.success].items = append(monkeys[m.success].items, item)
				} else {
					monkeys[m.failure].items = append(monkeys[m.failure].items, item)
				}
			}

			monkeys[i].items = nil
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(inspected)))
	return inspected[0] * inspected[1]
}

func part2(monkeys []monkey) int {
	lcm := 1
	for _, m := range monkeys {
		lcm *= m.test
	}

	var inspected []int
	for range monkeys {
		inspected = append(inspected, 0)
	}

	for round := 0; round < 10_000; round++ {
		for i, m := range monkeys {
			for _, item := range m.items {
				inspected[i]++

				var a int
				if m.operands[0].old {
					a = item
				} else {
					a = m.operands[0].value
				}

				var b int
				if m.operands[1].old {
					b = item
				} else {
					b = m.operands[1].value
				}

				switch m.operation {
				case "+":
					item = a + b
				case "*":
					item = a * b
				}

				if item%m.test == 0 {
					monkeys[m.success].items = append(monkeys[m.success].items, item%lcm)
				} else {
					monkeys[m.failure].items = append(monkeys[m.failure].items, item%lcm)
				}
			}

			monkeys[i].items = nil
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(inspected)))
	return inspected[0] * inspected[1]
}

func main() {
	monkeys := parse()
	fmt.Println(part1(monkeys))
	monkeys = parse()
	fmt.Println(part2(monkeys))
}
