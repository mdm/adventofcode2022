package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type fileinfo struct {
	size int
	name string
}

type directory struct {
	files   []fileinfo
	subdirs []string
}

func parse() (fs map[string]directory) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fs = make(map[string]directory)

	var cwd []string
	var files []fileinfo
	var subdirs []string
	listing := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		if listing {
			switch parts[0] {
			case "$":
				path := strings.Join(cwd, "")
				dir := directory{files, subdirs}
				fs[path] = dir
				files = make([]fileinfo, 0)
				subdirs = make([]string, 0)
				listing = false
			case "dir":
				subdirs = append(subdirs, parts[1]+"/")
			default:
				size, err := strconv.Atoi(parts[0])
				if err != nil {
					log.Fatal(err)
				}
				file := fileinfo{size, parts[1]}
				files = append(files, file)
			}
		}

		if !listing {
			switch parts[1] {
			case "cd":
				switch parts[2] {
				case "/":
					cwd = []string{"/"}
				case "..":
					cwd = cwd[:len(cwd)-1]
				default:
					cwd = append(cwd, parts[2]+"/")
				}
			case "ls":
				listing = true
			}
		}
	}

	if listing {
		path := strings.Join(cwd, "")
		dir := directory{files, subdirs}
		fs[path] = dir
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return fs
}

func totalSize(path string, dir directory, fs map[string]directory) int {
	size := 0

	for _, file := range dir.files {
		size += file.size
	}

	for _, subdir := range dir.subdirs {
		subpath := path + subdir
		size += totalSize(subpath, fs[subpath], fs)
	}

	return size
}

func part1(fs map[string]directory) int {
	answer := 0
	for path, dir := range fs {
		size := totalSize(path, dir, fs)
		if size <= 100_000 {
			answer += size
		}
	}

	return answer
}

func part2(fs map[string]directory) int {
	min := totalSize("/", fs["/"], fs)
	unused := 70_000_000 - min
	needed := 30_000_000 - unused

	for path, dir := range fs {
		size := totalSize(path, dir, fs)
		if size < min && size > needed {
			min = size
		}
	}

	return min
}

func main() {
	fs := parse()
	fmt.Println(part1(fs))
	fmt.Println(part2(fs))
}
