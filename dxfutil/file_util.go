package dxfutil


import (
	"os"
	"bufio"
	"strings"
)

func PathExists(path string) (bool) {
	_, err := os.Stat(path)
	r := false
	if err == nil {r = true}
	if os.IsNotExist(err) {r = false}
	return r
}

func Readini(m map[string]string, path string) {
	path += "/dxf2gcode.ini"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
		os.Exit(0)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		for range parts {
			m[parts[0]] = parts[1]
		}
	}
}

