package main

import (
	"os"
	"os/user"
	"fmt"
	"github.com/jethornton/dxfutil"
)

func main(){
	usr, _ := user.Current()
	iniMap := make(map[string]string)
	var inFile string
	if len(os.Args) == 2 {
		switch os.Args[1]{
		case "-v":
			fmt.Println("Version 0.001")
			os.Exit(0)
		default:
			inFile = os.Args[1]
		}
	} else {
		pwd, _ := os.Getwd()
		fmt.Println("Current Working Directory is:", pwd)
		fmt.Println("Current User Directory is:", usr.HomeDir)
		fmt.Println("Usage is: dxf2gcode filename.ext")
		fmt.Println("Usage is: dxf2gcode -v")
		os.Exit(0)
	}
	dxfutil.Readini(iniMap, usr.HomeDir)
	lines := dxfutil.GetLines(inFile)
	entities := dxfutil.GetEntities(lines)
	entities = dxfutil.GetEndPoints(entities)
	start := 1
	entities = dxfutil.GetIndex(start, iniMap["TOLERANCE"], entities)
	dxfutil.GenGcode(entities, iniMap["SAVEAS"])
}
