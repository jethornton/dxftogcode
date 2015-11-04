package dxfutil

import (
	"bufio"
	"os"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sort"
)

type Ent struct {
	Test, Index int
	G0, G, G8, G10, G11, G20, G21, G30, G31, G40, G50, G51 string
	Xs, Xe, Ys, Ye, Zs, Ze float64
}

type ByIndex []Ent

func (a ByIndex) Len() int { return len(a) }
func (a ByIndex) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

// add float compare 

var EPSILON float64 = 0.0001

func floatCompare(a, b float64) bool {
	if ((a - b) < EPSILON && (b - a) < EPSILON) {
		return true
	}
	return false
} 

func formatString (s string) (r string) {
	p := strings.Index(s, ".")
	if p == -1 { 
		s += "."
		p = strings.Index(s, ".")
	}
	for i := len(s[p+1:])-5 ; i < 0 ; i++ {
		s += "0"
	}
	if len(s[p+1:]) > 5 {s = s[:p+6]} // trim decimal places to 5
	return s
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) (float64) {
	shift := math.Pow(10, float64(places))
	return Round(f * shift) / shift;
}

func getXcartesian(r, a, x string) (float64){
	radius, _ := strconv.ParseFloat(r, 64)
	angle, _ := strconv.ParseFloat(a, 64)
	offset, _ := strconv.ParseFloat(x, 64)
	return RoundPlus(radius * math.Cos(angle * (math.Pi / 180)), 4) + offset
}
func getYcartesian(r, a, y string) (c float64){
	radius, _ := strconv.ParseFloat(r, 64)
	angle, _ := strconv.ParseFloat(a, 64)
	offset, _ := strconv.ParseFloat(y, 64)
return RoundPlus(radius * math.Sin(angle * (math.Pi / 180)), 4) + offset
}

func GetLines(f string) ([]string) {
	lines := []string{}
	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	inEntities := false
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
	fmt.Println("GetLines Done")
	return lines
}

func GetEntities(list []string) ([]Ent){
	// add error handling
	group := ""
	count := 0
	var e []Ent
	for i := range list {
		if len(group) > 0 {
			switch group {
			case "  0":
				if len(e) > 0 {count++}
				e = append(e, Ent{})
				e[count].Test = count
				e[count].G0 = list[i]
			case "  8":
				e[count].G8 = list[i]
			case " 10":
				e[count].G10 = formatString(list[i])
			case " 11":
				e[count].G11 = formatString(list[i])
			case " 20":
				e[count].G20 = formatString(list[i])
			case " 21":
				e[count].G21 = formatString(list[i])
			case " 30":
				e[count].G30 = formatString(list[i])
			case " 31":
				e[count].G31 = formatString(list[i])
			case " 40":
				e[count].G40 = formatString(list[i])
			case " 50":
				e[count].G50 = formatString(list[i])
			case " 51":
				e[count].G51 = formatString(list[i])
			}
		group = ""
		}
		switch list[i] { // trigger when a group is found
		case "  0", "  8", " 10", " 11", " 20", " 21",
		 " 30", " 31", " 40", " 50", " 51":
			group = list[i]
		}
	}
	fmt.Println("GetEntities Done")
	return e
}

func GetLayers(e []Ent){
	var layers []string
	Search:
	for i := range e {
		if len(layers) == 0 {
			layers = append(layers, e[i].G8)
			continue Search
		}
		for j := range layers {
			if e[i].G8 == layers[j]{ continue Search }
		}
		layers = append(layers, e[i].G8)
	}
	fmt.Println(layers)
}

func GetEndPoints (e []Ent) ([]Ent){
	// add error handling
	for i := range e {
		switch e[i].G0 {
		case "ARC": // get the X and Y end points
			e[i].Xs = getXcartesian(e[i].G40, e[i].G50, e[i].G10)
			e[i].Xe = getXcartesian(e[i].G40, e[i].G51, e[i].G10)
			e[i].Ys = getYcartesian(e[i].G40, e[i].G50, e[i].G20)
			e[i].Ye = getYcartesian(e[i].G40, e[i].G51, e[i].G20)
		case "LINE":
			e[i].Xs, _ = strconv.ParseFloat(e[i].G10, 64)
			e[i].Ys, _ = strconv.ParseFloat(e[i].G20, 64)
			e[i].Zs, _ = strconv.ParseFloat(e[i].G30, 64)
			e[i].Xe, _ = strconv.ParseFloat(e[i].G11, 64)
			e[i].Ye, _ = strconv.ParseFloat(e[i].G21, 64)
			e[i].Zs, _ = strconv.ParseFloat(e[i].G31, 64)
		case "CIRCLE":
			fmt.Println("Circles not supported at this time")
			os.Exit(1)
		case "SPLINE":
			fmt.Println("Splines not supported at this time")
			os.Exit(1)
		}
	}
	fmt.Println("GetEndPoints Done")
	return e
}

func GetOrder(e []Ent) ([]Ent) {
	c := 3 // entity to search from
	dir := "CW" // direction of travel
	Search:
	for i := range e {
		if i == len(e) - 1 { break } // don't process the last one
		// if direction is CW and it is an arc reverse the arc before processing
		if i == 0 { // this will need to be smarter to figure out if it is G2 or G3
			switch e[c].G0 {
			case "ARC":
				if dir == "CW" {
					e[c].G = "2"
					e[c].Xe, e[c].Xs = e[c].Xs, e[c].Xe
					e[c].Ye, e[c].Ys = e[c].Ys, e[c].Ye
					e[c].G50, e[c].G51 = e[c].G51, e[c].G50
				} else {
					e[c].G = "3"
				}
			case "LINE":
				e[c].G = "1"
			}
		}
		for j := range e {
			if c != j && floatCompare(e[c].Xe, e[j].Xs) && floatCompare(e[c].Ye, e[j].Ys) {
				e[j].Index = i + 1
				c = j
				switch e[j].G0 {
				case "ARC":
					e[j].G = "3"
				case "LINE":
					e[j].G = "1"
				}
				continue Search
			}
		}
		for k := range e {
			if c != k && floatCompare(e[c].Xe, e[k].Xe) && floatCompare(e[c].Ye, e[k].Ye) {
				// swap end points
				e[k].Xe, e[k].Xs = e[k].Xs, e[k].Xe
				e[k].Ye, e[k].Ys = e[k].Ys, e[k].Ye
				
				e[k].Index = i + 1
				c = k
				switch e[k].G0 {
				case "ARC":
					e[k].G = "2"
					e[k].G50, e[k].G51 = e[k].G51, e[k].G50
				case "LINE":
					e[k].G = "1"
					e[k].G10, e[k].G11 = e[k].G11, e[k].G10
					e[k].G20, e[k].G21 = e[k].G21, e[k].G20
				}
				continue Search
			}
		}
	}
	sort.Sort(ByIndex(e))
	fmt.Println("GetOrder Done")
	return e
}
