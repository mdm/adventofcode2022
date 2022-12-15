package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type position struct {
	x int
	y int
}

type sensor struct {
	location position
	beacon   position
}

func parse() (sensors []sensor) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	regex := regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		match := regex.FindStringSubmatch(line)

		lx, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}

		ly, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}

		bx, err := strconv.Atoi(match[3])
		if err != nil {
			log.Fatal(err)
		}

		by, err := strconv.Atoi(match[4])
		if err != nil {
			log.Fatal(err)
		}

		sensors = append(sensors, sensor{position{lx, ly}, position{bx, by}})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return sensors
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

type interval struct {
	min int
	max int
}

func countCovered(sensors []sensor, row int) int {
	var projections []interval
	for _, sensor := range sensors {
		sensorRange := absDiffInt(sensor.location.x, sensor.beacon.x) + absDiffInt(sensor.location.y, sensor.beacon.y)
		rowDistance := absDiffInt(sensor.location.y, row)
		overshoot := sensorRange - rowDistance

		if overshoot >= 0 {
			projection := interval{sensor.location.x - overshoot, sensor.location.x + overshoot}
			projections = append(projections, projection)
		}
	}

	var oldChecked []interval
	for _, projection := range projections {
		var newChecked []interval
		skip := false
		for _, checked := range oldChecked {
			if projection.min >= checked.min && projection.max <= checked.max {
				// projection completely contained in checked
				skip = true
				break
			}

			if projection.min > checked.max || projection.max < checked.min {
				// completely disjunct
				newChecked = append(newChecked, checked)
				continue
			}

			if checked.min < projection.min {
				newChecked = append(newChecked, interval{checked.min, projection.min - 1})
			}

			if checked.max > projection.max {
				newChecked = append(newChecked, interval{projection.max + 1, checked.max})
			}
		}

		if !skip {
			newChecked = append(newChecked, projection)
			oldChecked = newChecked
		}
	}

	covered := 0
	for _, checked := range oldChecked {
		covered += checked.max - checked.min + 1
	}

	return covered
}

func part1(sensors []sensor) int {
	var row int
	if os.Args[1] == "input.txt" {
		row = 2_000_000
	} else {
		row = 10
	}

	covered := countCovered(sensors, row)

	onRow := make(map[position]bool)
	for _, sensor := range sensors {
		if sensor.beacon.y == row {
			onRow[sensor.beacon] = true
		}
	}

	return covered - len(onRow)
}

func findDistressBeacon(sensors []sensor, row int, window interval) (int, []interval) {
	var projections []interval
	for _, sensor := range sensors {
		sensorRange := absDiffInt(sensor.location.x, sensor.beacon.x) + absDiffInt(sensor.location.y, sensor.beacon.y)
		rowDistance := absDiffInt(sensor.location.y, row)
		overshoot := sensorRange - rowDistance

		if overshoot >= 0 {
			min := window.min
			if sensor.location.x-overshoot > min {
				min = sensor.location.x - overshoot
			}

			max := window.max
			if sensor.location.x+overshoot < max {
				max = sensor.location.x + overshoot
			}

			projection := interval{min, max}
			projections = append(projections, projection)
		}
	}

	var oldChecked []interval
	for _, projection := range projections {
		var newChecked []interval
		skip := false
		for _, checked := range oldChecked {
			if projection.min >= checked.min && projection.max <= checked.max {
				// projection completely contained in checked
				skip = true
				break
			}

			if projection.min > checked.max || projection.max < checked.min {
				// completely disjunct
				newChecked = append(newChecked, checked)
				continue
			}

			if checked.min < projection.min {
				newChecked = append(newChecked, interval{checked.min, projection.min - 1})
			}

			if checked.max > projection.max {
				newChecked = append(newChecked, interval{projection.max + 1, checked.max})
			}
		}

		if !skip {
			newChecked = append(newChecked, projection)
			oldChecked = newChecked
		}
	}

	covered := 0
	for _, checked := range oldChecked {
		covered += checked.max - checked.min + 1
	}

	return covered, oldChecked
}

func part2(sensors []sensor) int {
	var window int
	if os.Args[1] == "input.txt" {
		window = 4_000_000
	} else {
		window = 20
	}

	for y := 0; y <= window; y++ {
		covered, projections := findDistressBeacon(sensors, y, interval{0, window})
		if covered < window+1 {
			for x := 0; x <= window; x++ {
				found := true
				for _, projection := range projections {
					if x >= projection.min && x <= projection.max {
						found = false
						break
					}
				}

				if found {
					return x*4_000_000 + y
				}
			}
		}
	}

	return -1
}

func main() {
	sensors := parse()
	fmt.Println(part1(sensors))
	fmt.Println(part2(sensors))
}
