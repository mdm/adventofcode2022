package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func parse() []string {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func snafu2decimal(snafu string) int {
	digits := "=-012"

	decimal := 0
	for _, digit := range snafu {
		decimal *= 5
		decimal += strings.Index(digits, string(digit)) - 2
	}

	return decimal
}

func reverse(str string) string {
	var result string
	for _, v := range str {
		result = string(v) + result
	}
	return result
}

func decimal2snafu(decimal int) string {
	digits := "=-012"

	var base5 []int
	for decimal > 0 {
		base5 = append(base5, decimal%5)
		decimal /= 5
	}

	snafu := ""
	carry := 0
	for _, digit := range base5 {
		digit += carry % 5
		carry /= 5
		if digit < 3 {
			snafu += string(digits[digit+2])
		} else {
			snafu += string(digits[digit-5+2])
			carry++
		}
	}

	for carry > 0 {
		digit := carry % 5
		carry /= 5
		if digit < 3 {
			snafu += string(digits[digit+2])
		} else {
			snafu += string(digits[digit-5+2])
			carry++
		}
	}

	return reverse(snafu)
}

func part1(snafus []string) string {
	sum := 0
	for _, snafu := range snafus {
		sum += snafu2decimal(snafu)
	}
	return decimal2snafu(sum)
}

func main() {
	snafus := parse()
	fmt.Println(part1(snafus))
}
