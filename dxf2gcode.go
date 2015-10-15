package main
// dxf2g version 0.0000001
import (
	"bufio"
	"os"
	"fmt"
	"strconv"
	"math"
	"sort"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
		os.Exit(0)
	}
}

func iniRead(m map[string]string) {
	f, err := os.Open("dxf2gcode.ini")
	check(err)
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

func getXcartesian(radius, angle float64) (coordinate float64){
	return radius * math.Cos(angle * (math.Pi / 180))
}
func getYcartesian(radius, angle float64) (coordinate float64){
	return radius * math.Sin(angle * (math.Pi / 180))
}

func getLines(filename string) ([]string) {
	lines := []string{}
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	inEntities := false
	 // extract all the lines of entities
	for scanner.Scan() {
		switch scanner.Text() {
		case "ENTITIES":
			inEntities = true
		case "ENDSEC":
			inEntities = false
		}
		if inEntities {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}

func getEnt(list []string, e []Ent) ([]Ent){
	// add error handling
	group := ""
	count := 0
	for i := range list {
		if len(group) > 0 {
			switch group {
			case "  0":
				if len(e) > 0 {count++}
				e = append(e, Ent{})
				e[count].G0 = list[i]
			case "  8":
				e[count].G8 = list[i]
			case " 10":
				e[count].G10, _ = strconv.ParseFloat(list[i],64)
			case " 11":
				e[count].G11, _ = strconv.ParseFloat(list[i],64)
			case " 20":
				e[count].G20, _ = strconv.ParseFloat(list[i],64)
			case " 21":
				e[count].G21, _ = strconv.ParseFloat(list[i],64)
			case " 30":
				e[count].G30, _ = strconv.ParseFloat(list[i],64)
			case " 31":
				e[count].G31, _ = strconv.ParseFloat(list[i],64)
			case " 40":
				e[count].G40, _ = strconv.ParseFloat(list[i],64)
			case " 50":
				e[count].G50, _ = strconv.ParseFloat(list[i],64)
			case " 51":
				e[count].G51, _ = strconv.ParseFloat(list[i],64)
			}
		group = ""
		}
		switch list[i] { // trigger when a group is found
		case "  0", "  8", " 10", " 11", " 20", " 21",
		 " 30", " 31", " 40", " 50", " 51":
			group = list[i]
		}
	}
	return e
}

func findEndPoints (e []Ent) ([]Ent){
	// add error handling
	for i := range e {
		switch e[i].G0 {
		case "ARC": // get the X and Y end points
			e[i].Xs = getXcartesian(e[i].G40, e[i].G50) + e[i].G10
			e[i].Xe = getXcartesian(e[i].G40, e[i].G51) + e[i].G10
			e[i].Ys = getYcartesian(e[i].G40, e[i].G50) + e[i].G20
			e[i].Ye = getYcartesian(e[i].G40, e[i].G51) + e[i].G20
		case "LINE": 
			e[i].Xs = e[i].G10
			e[i].Ys = e[i].G20
			e[i].Zs = e[i].G30
			e[i].Xe = e[i].G11
			e[i].Ye = e[i].G21
			e[i].Zs = e[i].G31
		case "CIRCLE": // if it is a circle it must be the only entity on that layer
			fmt.Println("Circle")
		}
	}
	return e
}

// find the matching start points for each entity and assign index
func findIndex(s int, t string, e []Ent) ([]Ent) {
	// add error handling
	xe := e[s].Xe
	ye := e[s].Ye
	index := 1
	tol, err := strconv.ParseFloat(t, 64)
	check(err)
	for skip, _ := range e {
		if skip != s {
			for i := range e {
				if math.Abs(e[i].Xs - xe) < tol && math.Abs(e[i].Ys - ye) < tol {
					e[i].Index = index
					index++
					xe = e[i].Xe
					ye = e[i].Ye
					break
				}
			}
		}
	}
	return e
}

type Ent struct {
	Index int
	G0, G8 string
	G10, G11, G20, G21, G30, G31, G40, G50, G51,
	Xs, Xe, Ys, Ye, Zs, Ze float64
}

type ByIndex []Ent

func (a ByIndex) Len() int { return len(a) }
func (a ByIndex) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

func genGcode (e []Ent, fileout string) {
/*
need to figure out if the start point of first entity and the end point
of the last entity are the same if so then make the last move to the
start point.
*/
	xo, yo := 0.0, 0.0
	feed := 25.0
	//outfile := "/home/john/linuxcnc/nc_files/output.ngc"
	f, err := os.Create(fileout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintf(f, "G0 X%f Y%f\n", e[0].Xs, e[0].Ys)
	fmt.Fprintf(f, "F%.1f\n",feed)
	for _, ent := range e {
		switch ent.G0 {
		case "LINE":
			fmt.Fprintf(f, "G1 X%.4f Y%.4f\n", ent.G11, ent.G21)
		case "ARC":
			switch {
			case ent.G50 <= 90.0:
				xo = -(ent.Xs - ent.G10)
				yo = -(ent.Ys - ent.G20)
			case ent.G50 > 90 && ent.G50 <= 180.0:
				xo = math.Abs(ent.Xs - ent.G10)
				yo = -(ent.Ys - ent.G20)
			case ent.G50 > 180 && ent.G50 <= 270.0:
				xo = math.Abs(ent.Xs - ent.G10)
				yo = math.Abs(ent.Ys - ent.G20)
			case ent.G50 > 270 && ent.G50 <= 360.0:
				xo = -(ent.Xs - ent.G10)
				yo = -(ent.Ys - ent.G20)
			}
			fmt.Fprintf(f, "G3 X%.4f Y%.4f I%.4f J%.4f\n", ent.Xe, ent.Ye, xo, yo)
		case "CIRCLE":
			fmt.Println("Circle")
		}
	}
	fmt.Fprintf(f, "M2")
	fmt.Println("Processing Done.")
}

func main(){
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
		//inFile = "test.dxf"
		fmt.Println("Usage is: dxf2g filename.ext")
		fmt.Println("Usage is: dxf2g -v")
		os.Exit(0)
	}
	iniRead(iniMap)
	fmt.Println(iniMap)
	fmt.Println(iniMap["SAVEAS"])
	var entities []Ent
	lines := getLines(inFile)
	entities = getEnt(lines, entities)
	entities = findEndPoints(entities)
	start := 1
	entities = findIndex(start, iniMap["TOLERANCE"], entities)
	sort.Sort(ByIndex(entities))
	genGcode(entities, iniMap["SAVEAS"])
}
