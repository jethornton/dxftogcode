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
	f, err := os.Open(path)
	if err != nil {
		panic(err)
		os.Exit(0)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "[Configuration]" && len(line) > 0 {
			parts := strings.Split(line, "=")
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
}

