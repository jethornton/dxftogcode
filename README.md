# DXFtoGcode
Don't confuse this with other dxf2gcode programs.

Status: Not Functioning

Current: adding command line flags and analyze

At this time dxf2gcode will convert a single layer DXF file to G code.

This is for G17 applications only at this time.

Usage:
GUI dxf.py

Comand Line Flags
Usage of dxf2gcode:
  -a	Analyze contents of the file
  -c	Convert contents of the file
  -d string
    	Direction of path (default "CCW")
  -i string
    	Input file path (default "/dxf/test.dxf")
  -o string
    	Output file path (default "output.ngc")

