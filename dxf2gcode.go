package main

// DXF to G code converter
import (
	"flag"
	"path/filepath"
	"log"
	"fmt"
	"os"
	"os/user"
	"github.com/jethornton/dxf2gcode/dxfutil"
)

var ( // flag variables
	input *string
	direction *string
	output *string
	analyze *bool
	convert *bool
	//port *int
)

// Basic flag declarations are available for string, integer, and boolean options.
func init() { // flag.Type(flag, default, help string)
	input = flag.String("i", "dxf/test.dxf", "Input file path")
	output = flag.String("o", "output.ngc", "Output file path")
	direction = flag.String("d", "CCW", "Direction of path")
	analyze = flag.Bool("a", false, "Analyze contents of the file")
	convert = flag.Bool("c", false, "Convert contents of the file")
	//port = flag.Int("port", 3000, "an int")
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 { // if no flags are passed print usage
		flag.Usage()
		fmt.Println("Analyze", *analyze)
		fmt.Println("Flags", flag.NFlag())
		os.Exit(1)
	}
	if *convert {
		fmt.Println("Convert was true")
	}
	iniMap := make(map[string]string)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dxfutil.Readini(iniMap, dir)
	lines := dxfutil.GetLines(*input)
	entities := dxfutil.GetEntities(lines)
	if *analyze {
		dxfutil.GetLayers(entities)
	}

	os.Exit(1)
	usr, _ := user.Current() // get user information
	inipath := usr.HomeDir + "/.config/dxf2emc"
	fmt.Println(inipath)
	//fmt.Println(dxfutil.PathExists(inipath + "/dxf2emc.ini"))
	//fmt.Println(dir)
	//cwd, _ := os.Getwd() // get current working directory

	//var inFile string
	/*
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-v":
			fmt.Println("Version 0.001")
			os.Exit(0)
		case "-p":
			fmt.Println("-p")
		default:
			inFile = os.Args[1]
		}
	} else {
		fmt.Println("Current Working Directory is:", cwd)
		fmt.Println("Current User Directory is:", usr.HomeDir)
		fmt.Println("Usage is: dxf2gcode filename.ext")
		fmt.Println("Usage is: dxf2gcode -v")
		os.Exit(0)
	}*/

	entities = dxfutil.GetEndPoints(entities)
	entities = dxfutil.GetOrder(entities)
	dxfutil.GenGcode(entities, iniMap["SAVEAS"])
/*
	for _, e := range entities {
		fmt.Printf("%2d %2d %4s Xs %9f Xe %9f Ys %9f Ye %9f\n",
		e.Test, e.Index, e.G0, e.Xs, e.Xe, e.Ys, e.Ye)
	}
*/
}
