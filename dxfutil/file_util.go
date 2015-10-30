package dxfutil


import (
	"os"
	"bufio"
	"strings"
)

func Readini(m map[string]string, home string) {
	home += "/dxf2gcode.ini"
	f, err := os.Open(home)
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

