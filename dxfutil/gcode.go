package dxfutil

import (
	"math"
	"fmt"
	"os"
	"strconv"
)

func GenGcode (e []Ent, fileout string) {
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

	fmt.Fprintf(f, "G0 X%.5f Y%.5f\n", e[0].Xs, e[0].Ys)
	fmt.Fprintf(f, "F%.1f\n",feed)
	for _, ent := range e {
		switch ent.G0 {
		case "LINE":
			fmt.Fprintf(f, "G%s X%s Y%s\n", ent.G, ent.G11, ent.G21)
		case "ARC":
			g10, _ := strconv.ParseFloat(ent.G10, 64)
			g20, _ := strconv.ParseFloat(ent.G20, 64)
			g50, _ := strconv.ParseFloat(ent.G50, 64)
			switch {
			case g50 <= 90.0:
				xo = -(ent.Xs - g10)
				yo = -(ent.Ys - g20)
			case g50 > 90 && g50 <= 180.0:
				xo = math.Abs(ent.Xs - g10)
				yo = -(ent.Ys - g20)
			case g50 > 180 && g50 <= 270.0:
				xo = math.Abs(ent.Xs - g10)
				yo = math.Abs(ent.Ys - g20)
			case g50 > 270 && g50 <= 360.0:
				xo = -(ent.Xs - g10)
				yo = -(ent.Ys - g20)
			}
			fmt.Fprintf(f, "G%s X%.5f Y%.5f I%.5f J%.5f\n", ent.G, ent.Xe, ent.Ye, xo, yo)
		case "CIRCLE":
			fmt.Println("Circle")
		}
	}
	fmt.Fprintf(f, "M2")
	fmt.Println("Processing Done.")
}
