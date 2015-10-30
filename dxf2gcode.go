package main

// DXF to G code converter
import (
	"fmt"
	"github.com/jethornton/dxf2gcode/dxfutil"
	"os"
	"os/user"
)

func main() {
	usr, _ := user.Current() // get user information
	cwd, _ := os.Getwd() // get current working directory
	iniMap := make(map[string]string)
	var inFile string
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-v":
			fmt.Println("Version 0.001")
			os.Exit(0)
		default:
			inFile = os.Args[1]
		}
	} else {
		fmt.Println("Current Working Directory is:", cwd)
		fmt.Println("Current User Directory is:", usr.HomeDir)
		fmt.Println("Usage is: dxf2gcode filename.ext")
		fmt.Println("Usage is: dxf2gcode -v")
		os.Exit(0)
	}
	dxfutil.Readini(iniMap, cwd)
	lines := dxfutil.GetLines(inFile)
	entities := dxfutil.GetEntities(lines)
	for _, e := range entities {
		fmt.Printf("%2d %4s G10 %8s G11 %8s G20 %8s G21 %8s G50 %9s G51 %9s\n",
		e.N, e.G0, e.G10, e.G11, e.G20, e.G21, e.G50, e.G51)
	}
	entities = dxfutil.GetEndPoints(entities)

	for _, e := range entities {
		fmt.Printf("%2d %4s Xs %9f Xe %9f Ys %9f Ye %9f\n",
		e.N, e.G0, e.Xs, e.Xe, e.Ys, e.Ye)
	}

	entities = dxfutil.GetIndex(entities)

	for _, e := range entities {
		fmt.Printf("%2d %2d %4s Xs %9f Xe %9f Ys %9f Ye %9f\n",
		e.N, e.Index, e.G0, e.Xs, e.Xe, e.Ys, e.Ye)
	}

	dxfutil.GenGcode(entities, iniMap["SAVEAS"])

}
