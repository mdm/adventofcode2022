package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
)

type element struct {
	singleton bool
	list      []element
	integer   int
}

func parseList(input string) (output element, rest string) {
	output = element{false, nil, 0}
	var acc []rune
	i := 0
	for i < len(input) {
		switch input[i] {
		case '[':
			i++
			innerOutput, innerRest := parseList(input[i:])
			output.list = append(output.list, innerOutput)
			i = 0
			input = innerRest
		case ']':
			i++
			if len(acc) > 0 {
				integer, err := strconv.Atoi(string(acc))
				if err != nil {
					log.Fatal(err)
				}

				output.list = append(output.list, element{true, nil, integer})
			}
			rest = input[i:]
			return output, rest
		case ',':
			i++
			if len(acc) > 0 {
				integer, err := strconv.Atoi(string(acc))
				if err != nil {
					log.Fatal(err)
				}

				output.list = append(output.list, element{true, nil, integer})
			}
			acc = nil
		default:
			acc = append(acc, rune(input[i]))
			i++
		}
	}

	rest = ""
	return output, rest
}

func parse() (leftPackets []element, rightPackets []element) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		line := scanner.Text()

		switch i % 3 {
		case 0:
			parsed, _ := parseList(line[1:])
			leftPackets = append(leftPackets, parsed)
		case 1:
			parsed, _ := parseList(line[1:])
			rightPackets = append(rightPackets, parsed)
		case 2:
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return leftPackets, rightPackets
}

func isOrdered(left element, right element) int {
	if left.singleton && right.singleton {
		if left.integer < right.integer {
			return 1
		}

		if left.integer > right.integer {
			return -1
		}

		return 0
	}

	if !left.singleton && !right.singleton {
		if len(left.list) < len(right.list) {
			for i, innerLeft := range left.list {
				innerRight := right.list[i]

				ordered := isOrdered(innerLeft, innerRight)
				if ordered != 0 {
					return ordered
				}
			}

			return 1
		}

		if len(left.list) > len(right.list) {
			for i, innerRight := range right.list {
				innerLeft := left.list[i]

				ordered := isOrdered(innerLeft, innerRight)
				if ordered != 0 {
					return ordered
				}
			}

			return -1
		}

		for i, innerLeft := range left.list {
			innerRight := right.list[i]

			ordered := isOrdered(innerLeft, innerRight)
			if ordered != 0 {
				return ordered
			}
		}

		return 0
	}

	if left.singleton {
		return isOrdered(element{false, []element{left}, 0}, right)
	} else {
		return isOrdered(left, element{false, []element{right}, 0})
	}
}

func part1(leftPackets []element, rightPackets []element) int {
	sum := 0
	for i, left := range leftPackets {
		right := rightPackets[i]

		if isOrdered(left, right) == 1 {
			sum += i + 1
		}
	}

	return sum
}

type signal []element

func (s signal) Len() int {
	return len(s)
}

func (s signal) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s signal) Less(i, j int) bool {
	return isOrdered(s[i], s[j]) == 1
}

func part2(leftPackets []element, rightPackets []element) int {
	dividers := []element{{false, []element{{true, nil, 2}}, 0}, {false, []element{{true, nil, 6}}, 0}}

	packets := append(leftPackets, rightPackets...)
	packets = append(packets, dividers...)

	sort.Sort(signal(packets))

	decoder := 1
	for i, packet := range packets {
		if reflect.DeepEqual(packet, dividers[0]) || reflect.DeepEqual(packet, dividers[1]) {
			decoder *= i + 1
		}
	}

	return decoder
}

func main() {
	leftPackets, rightPackets := parse()
	fmt.Println(part1(leftPackets, rightPackets))
	fmt.Println(part2(leftPackets, rightPackets))
}
