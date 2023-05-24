package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type vec2 struct {
	x int
	y int
}

type player struct {
	position    vec2
	orientation int
}

func parse() (board []string, size int, human player, instructions []string) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	size, err = strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	board = []string{}
	maxLen := 0
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			break
		}

		board = append(board, line)
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	instructions = []string{}
	if scanner.Scan() {
		line := scanner.Text()

		start := 0
		end := 0
		for end < len(line) {
			if line[end] == 'L' || line[end] == 'R' {
				if start < end {
					instructions = append(instructions, line[start:end])
				}

				instructions = append(instructions, string(line[end]))

				end++
				start = end
			} else {
				end++
			}
		}

		if start < end {
			instructions = append(instructions, line[start:end])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i, line := range board {
		padding := maxLen - len(line)
		board[i] += strings.Repeat(" ", padding)
	}

	human = player{vec2{0, 0}, 0}
	for x, tile := range board[0] {
		if tile == '.' {
			human.position.x = x
			break
		}
	}

	return board, size, human, instructions
}

type vec3 struct {
	x int
	y int
	z int
}

func (v *vec3) dot_product(other vec3) int {
	return v.x*other.x + v.y*other.y + v.z*other.z
}

type mat33 struct {
	rows []vec3
}

func (m *mat33) multiply_vec3(v vec3) vec3 {
	return vec3{
		m.rows[0].dot_product(v),
		m.rows[1].dot_product(v),
		m.rows[2].dot_product(v),
	}
}

func (m *mat33) multiply_mat33(other mat33) mat33 {
	return mat33{
		[]vec3{
			{
				m.rows[0].dot_product(vec3{other.rows[0].x, other.rows[1].x, other.rows[2].x}),
				m.rows[0].dot_product(vec3{other.rows[0].y, other.rows[1].y, other.rows[2].y}),
				m.rows[0].dot_product(vec3{other.rows[0].z, other.rows[1].z, other.rows[2].z}),
			},
			{
				m.rows[1].dot_product(vec3{other.rows[0].x, other.rows[1].x, other.rows[2].x}),
				m.rows[1].dot_product(vec3{other.rows[0].y, other.rows[1].y, other.rows[2].y}),
				m.rows[1].dot_product(vec3{other.rows[0].z, other.rows[1].z, other.rows[2].z}),
			},
			{
				m.rows[2].dot_product(vec3{other.rows[0].x, other.rows[1].x, other.rows[2].x}),
				m.rows[2].dot_product(vec3{other.rows[0].y, other.rows[1].y, other.rows[2].y}),
				m.rows[2].dot_product(vec3{other.rows[0].z, other.rows[1].z, other.rows[2].z}),
			},
		},
	}
}

func (m *mat33) transpose() mat33 {
	return mat33{
		[]vec3{
			{m.rows[0].x, m.rows[1].x, m.rows[2].x},
			{m.rows[0].y, m.rows[1].y, m.rows[2].y},
			{m.rows[0].z, m.rows[1].z, m.rows[2].z},
		},
	}
}

type tile struct {
	column    int
	row       int
	transform mat33
}

type face struct {
	center  vec3
	corners []vec3
}
type cube struct {
	faces     []vec3
	transform mat33
	tiles     []*tile
	turns     []rune // to rotate back to start orientation
}

func (c *cube) findFace(x int, y int, z int) int {
	for i, f := range c.faces {
		tf := c.transform.multiply_vec3(f)
		if tf.x == x && tf.y == y && tf.z == z {
			return i
		}
	}

	return -1
}

func (c *cube) rotateUp() cube {
	nextCube := *c

	rotation := mat33{
		[]vec3{
			{1, 0, 0},
			{0, 0, 1},
			{0, -1, 0},
		},
	}

	nextCube.transform = nextCube.transform.multiply_mat33(rotation)
	return nextCube
}

func (c *cube) rotateDown() cube {
	nextCube := *c

	rotation := mat33{
		[]vec3{
			{1, 0, 0},
			{0, 0, -1},
			{0, 1, 0},
		},
	}

	nextCube.transform = nextCube.transform.multiply_mat33(rotation)
	return nextCube
}

func (c *cube) rotateLeft() cube {
	nextCube := *c

	rotation := mat33{
		[]vec3{
			{0, 0, -1},
			{0, 1, 0},
			{1, 0, 0},
		},
	}

	nextCube.transform = nextCube.transform.multiply_mat33(rotation)
	return nextCube
}

func (c *cube) rotateRight() cube {
	nextCube := *c

	rotation := mat33{
		[]vec3{
			{0, 0, 1},
			{0, 1, 0},
			{-1, 0, 0},
		},
	}

	nextCube.transform = nextCube.transform.multiply_mat33(rotation)
	return nextCube
}

func makeCube(board []string, size int) cube {
	startCube := cube{
		[]vec3{
			// TODO: fix corners
			{0, 0, 1},  // 1
			{0, 1, 0},  // 2
			{-1, 0, 0}, // 3
			{1, 0, 0},  // 4
			{0, -1, 0}, // 5
			{0, 0, -1}, // 6
		},
		mat33{
			[]vec3{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
		[]*tile{nil, nil, nil, nil, nil, nil},
		nil,
	}

	startColumn := 0
	startRow := 0
FindStart:
	for y := 0; y < len(board)/size; y++ {
		for x := 0; x < len(board[0])/size; x++ {
			if board[y*size][x*size] != ' ' {
				startColumn = x
				startRow = y
				break FindStart
			}
		}
	}

	var transformRows []vec3
	transformRows = append(transformRows, startCube.transform.rows...)
	startCube.tiles[0] = &tile{
		startColumn,
		startRow,
		mat33{transformRows},
	}

	filled := 0
	var queue []cube
	queue = append(queue, startCube)
	for {
		current := queue[0]
		queue = queue[1:]

		frontFace := current.findFace(0, 0, 1)
		column := current.tiles[frontFace].column
		row := current.tiles[frontFace].row
		fmt.Println("*", frontFace, column, row, current.tiles)
		filled++

		if filled == 6 {
			return current
		}

		if (row+1) < len(board)/size && board[(row+1)*size][column*size] != ' ' {
			nextCube := current.rotateUp()
			nextFront := nextCube.findFace(0, 0, 1)
			if current.tiles[nextFront] == nil {
				var transformRows []vec3
				transformRows = append(transformRows, nextCube.transform.rows...)

				current.tiles[nextFront] = &tile{
					column,
					row + 1,
					mat33{transformRows},
				}
				queue = append(queue, nextCube)
			}
		}

		if (row-1) >= 0 && board[(row-1)*size][column*size] != ' ' {
			nextCube := current.rotateDown()
			nextFront := nextCube.findFace(0, 0, 1)
			if current.tiles[nextFront] == nil {
				var transformRows []vec3
				transformRows = append(transformRows, nextCube.transform.rows...)

				current.tiles[nextFront] = &tile{
					column,
					row - 1,
					mat33{transformRows},
				}
				queue = append(queue, nextCube)
			}
		}

		if (column+1) < len(board[0])/size && board[row*size][(column+1)*size] != ' ' {
			nextCube := current.rotateLeft()
			nextFront := nextCube.findFace(0, 0, 1)
			if current.tiles[nextFront] == nil {
				var transformRows []vec3
				transformRows = append(transformRows, nextCube.transform.rows...)

				current.tiles[nextFront] = &tile{
					column + 1,
					row,
					mat33{transformRows},
				}
				queue = append(queue, nextCube)
			}
		}

		if (column-1) >= 0 && board[row*size][(column-1)*size] != ' ' {
			nextCube := current.rotateRight()
			nextFront := nextCube.findFace(0, 0, 1)
			if current.tiles[nextFront] == nil {
				var transformRows []vec3
				transformRows = append(transformRows, nextCube.transform.rows...)

				current.tiles[nextFront] = &tile{
					column - 1,
					row,
					mat33{transformRows},
				}
				queue = append(queue, nextCube)
			}
		}
	}
}

func (p *player) move2d(board []string, distance int) {
	i := 0
	lastOpen := p.position
	for i < distance {
		done := false
		switch p.orientation {
		case 0:
			next := (p.position.x + 1) % len(board[p.position.y])
			if board[p.position.y][next] == '#' {
				done = true
				p.position = lastOpen
				break
			}

			p.position.x = next

			if board[p.position.y][next] == '.' {
				lastOpen.x = next
				i++
			}
		case 3:
			next := (p.position.y - 1 + len(board)) % len(board)
			if board[next][p.position.x] == '#' {
				done = true
				p.position = lastOpen
				break
			}

			p.position.y = next

			if board[next][p.position.x] == '.' {
				lastOpen.y = next
				i++
			}
		case 2:
			next := (p.position.x - 1 + len(board[p.position.y])) % len(board[p.position.y])
			if board[p.position.y][next] == '#' {
				done = true
				p.position = lastOpen
				break
			}

			p.position.x = next

			if board[p.position.y][next] == '.' {
				lastOpen.x = next
				i++
			}
		case 1:
			next := (p.position.y + 1) % len(board)
			if board[next][p.position.x] == '#' {
				done = true
				p.position = lastOpen
				break
			}

			p.position.y = next

			if board[next][p.position.x] == '.' {
				lastOpen.y = next
				i++
			}
		}

		if done {
			break
		}
	}
}

func (p *player) move3d(board []string, size int, cube cube, distance int) {
	i := 0
	for i < distance {
		break
	}
}

func (p *player) turn(direction int) {
	p.orientation = (p.orientation + direction + 4) % 4
}

func part1(board []string, human player, instructions []string) int {
	for _, instruction := range instructions {
		distance, err := strconv.Atoi(instruction)

		if err == nil {
			human.move2d(board, distance)
		} else {
			switch instruction {
			case "L":
				human.turn(-1)
			case "R":
				human.turn(1)
			}
		}
	}

	return (human.position.y+1)*1000 + (human.position.x+1)*4 + human.orientation
}

func part2(board []string, size int, cube cube, human player, instructions []string) int {
	fmt.Println()
	for i, t := range cube.tiles {
		fmt.Println(i, t.column, t.row, t.transform)
	}

	for _, instruction := range instructions {
		distance, err := strconv.Atoi(instruction)

		if err == nil {
			human.move3d(board, size, cube, distance)
		} else {
			switch instruction {
			case "L":
				human.turn(-1)
			case "R":
				human.turn(1)
			}
		}
	}

	return (human.position.y+1)*1000 + (human.position.x+1)*4 + human.orientation
}

func main() {
	board, size, human, instructions := parse()
	fmt.Println(part1(board, human, instructions))
	cube := makeCube(board, size)
	fmt.Println(part2(board, size, cube, human, instructions))
}
