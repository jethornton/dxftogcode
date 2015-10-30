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
	N, Index int
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
				e[count].N = count
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

func GetIndex(e []Ent) ([]Ent) {
	test := e
	m := 0
	Searching:
	for key := range e {
		for i := range test {
			if m != i {
				 // CCW match
				if floatCompare(e[m].Xe, e[i].Xs) && floatCompare(e[m].Ye, e[i].Ys) {
					e[i].Index = key
					m = i
					//fmt.Printf("CCW Match Xe%f Xs%f Ye%f Ys%f\n"
					//,e[m].Xe, e[i].Xs, e[m].Ye, e[i].Ys)
					switch e[m].G0 {
					case "ARC":
						e[m].G = "3"
					case "LINE":
						e[m].G = "1"
					}
					continue Searching
				}
			}
		}
		for i := range test {
			if m != i {
				 // CW match
				if floatCompare(e[m].Xe, e[i].Xe) && floatCompare(e[m].Ye, e[i].Ye) {
					e[i].Index = key
					m = i
					switch e[m].G0 {
					case "ARC":
						e[m].G = "2"
					case "LINE":
						e[m].G = "1"
					}
					// swap end points and angles
					e[m].Xe, e[m].Xs = e[m].Xs, e[m].Xe
					e[m].Ye, e[m].Ys = e[m].Ys, e[m].Ye
					e[m].G50, e[m].G51 = e[m].G51, e[m].G50
					//fmt.Printf("CW Match Xe%f Xs%f Ye%f Ys%f\n",e[m].Xe, e[i].Xs, e[m].Ye, e[i].Ys)
					continue Searching
				}
			}
		}
		fmt.Println("no match for",key)
	}
	sort.Sort(ByIndex(e))
	fmt.Println("GetIndex Done")
	return e
}
