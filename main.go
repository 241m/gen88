package main

import (
	"crypto/sha256"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type SvgDef struct {
	Id  string `xml:"id,attr"`
	Xml []byte `xml:",innerxml"`
}

type Svg struct {
	Defs []SvgDef `xml:"defs"`
}

const (
	size = 512      // Size (with and height) of resulting SVG.
	grid = size / 8 // Size of single symbol in grid.
)

const svgHeader = `<svg
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	width="%[1]d"
	height="%[1]d"
	fill="%s">
	<title>%x</title>
`

func main() {
	var symFile = flag.String("s", "symbols.svg", "Path of SVG doc with 16 symbol defs")
	var dataFile = flag.String("f", "-", `Input file (use "-" for stdin)`)
	var color = flag.String("c", "black", "Default fill color")

	flag.Parse()

	if sum, err := sha256File(*dataFile); err != nil {
		log.Fatal(err)
	} else if mat, err := createMatrix(sum); err != nil {
		log.Fatal(err)
	} else if err := writeSVG(mat, *symFile, *color, sum); err != nil {
		log.Fatal(err)
	}
}

// Get SHA256 sum of data from given file. If file path is
// "-", will get data from stdin instead.
func sha256File(path string) (sum *[32]byte, err error) {
	var data []byte

	if path == "-" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}

	if err == nil {
		arr := sha256.Sum256(data)
		sum = &arr
	}

	return
}

// Create an 8x8 matrix of int64 numbers from given SHA256 sum.
func createMatrix(sum *[32]byte) (*[8][8]int64, error) {
	mat := [8][8]int64{}

	for i, n := range sum {
		h := fmt.Sprintf("%02x", n)

		if err := addSymbol(&mat, h[0:1], i*2); err != nil {
			return nil, err
		}

		if err := addSymbol(&mat, h[1:], i*2+1); err != nil {
			return nil, err
		}
	}

	return &mat, nil
}

// Add symbol numered h (in hexadecimal) to the appropriate
// x, y position in matrix mat based on the position i in
// a flat 1D array.
func addSymbol(mat *[8][8]int64, h string, i int) error {
	// convert the hex back to an int64 (this gives us an
	// integer in the range of 0-15, this is used as the
	// symbol ID)
	if v, err := strconv.ParseInt(h, 16, 64); err != nil {
		return err
	} else {
		// add the symbol ID to the appropriate x, y
		// position in the matrix.
		x := i % 8
		y := i / 8
		mat[x][y] = v
	}
	return nil
}

// Write symbol matrix as an SVG to stdout.
func writeSVG(mat *[8][8]int64, sym string, fill string, sum *[32]byte) error {
	svgDefs := Svg{}

	if symSvg, err := os.ReadFile(sym); err != nil {
		return err
	} else if err := xml.Unmarshal(symSvg, &svgDefs); err != nil {
		return err
	}

	fmt.Printf(svgHeader, size, fill, *sum)

	for _, defs := range svgDefs.Defs {
		if len(defs.Id) > 0 {
			fmt.Printf(`<defs id="%s">%s</defs>`, defs.Id, defs.Xml)
		} else {
			fmt.Printf(`<defs>%s</defs>`, defs.Xml)
		}
		fmt.Println()
	}

	for y, row := range mat {
		for x, n := range row {
			fmt.Printf(
				`<use x="%d" y="%d" href="#sym-%[3]d" xlink:href="#sym-%[3]d" />%s`,
				x*grid, y*grid, n, "\n",
			)
		}
	}

	fmt.Println(`</svg>`)

	return nil
}
