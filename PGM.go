package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM represents a PGM image.
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	scanner.Scan()
	magicNumber := scanner.Text()

	// Ignore comments
	for strings.HasPrefix(scanner.Text(), "#") {
		scanner.Scan()
	}

	// Read width, height, and max value
	scanner.Scan()
	size := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(size[0])
	height, _ := strconv.Atoi(size[1])

	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())

	// Read image data
	var data [][]uint8
	if magicNumber == "P2" {
		data, _ = readP2(scanner, width, height)
	} else if magicNumber == "P5" {
		data, _ = readP5(file, width, height)
	} else {
		return nil, fmt.Errorf("unsupported PGM format: %s", magicNumber)
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         max,
	}, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write magic number, width, height, and max value
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// Write image data
	if pgm.magicNumber == "P2" {
		writeP2(writer, pgm.data)
	} else if pgm.magicNumber == "P5" {
		writeP5(file, pgm.data)
	} else {
		return fmt.Errorf("unsupported PGM format: %s", pgm.magicNumber)
	}

	writer.Flush()
	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pgm.data[y][x] = uint8(pgm.max) - pgm.data[y][x]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width/2; x++ {
			pgm.data[y][x], pgm.data[y][pgm.width-x-1] = pgm.data[y][pgm.width-x-1], pgm.data[y][x]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	for y := 0; y < pgm.height/2; y++ {
		pgm.data[y], pgm.data[pgm.height-y-1] = pgm.data[pgm.height-y-1], pgm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width)
	for x := 0; x < pgm.width; x++ {
		newData[x] = make([]uint8, pgm.height)
		for y := 0; y < pgm.height; y++ {
			newData[x][y] = pgm.data[pgm.height-y-1][x]
		}
	}
	pgm.data = newData
	pgm.width, pgm.height = pgm.height, pgm.width
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	data := make([][]bool, pgm.height)
	for y := 0; y < pgm.height; y++ {
		data[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			data[y][x] = pgm.data[y][x] > uint8(pgm.max/2)
		}
	}
	return &PBM{
		data:        data,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P4",
	}
}

// Additional helper functions

func readP2(scanner *bufio.Scanner, width, height int) ([][]uint8, error) {
	data := make([][]uint8, height)
	for y := 0; y < height; y++ {
		data[y] = make([]uint8, width)
		lineValues := strings.Fields(scanner.Text())
		for x := 0; x < width; x++ {
			value, err := strconv.Atoi(lineValues[x])
			if err != nil {
				return nil, err
			}
			data[y][x] = uint8(value)
		}
		scanner.Scan()
	}
	return data, nil
}

func readP5(file *os.File, width, height int) ([][]uint8, error) {
	data := make([][]uint8, height)
	for y := 0; y < height; y++ {
		data[y] = make([]uint8, width)
		buf := make([]byte, width)
		_, err := file.Read(buf)
		if err != nil {
			return nil, err
		}
		for x := 0; x < width; x++ {
			data[y][x] = uint8(buf[x])
		}
	}
	return data, nil
}

func writeP2(writer *bufio.Writer, data [][]uint8) {
	for y := 0; y < len(data); y++ {
		for x := 0; x < len(data[y]); x++ {
			fmt.Fprintf(writer, "%d ", data[y][x])
		}
		fmt.Fprintln(writer)
	}
}

func writeP5(file *os.File, data [][]uint8) {
	for y := 0; y < len(data); y++ {
		buf := make([]byte, len(data[y]))
		for x := 0; x < len(data[y]); x++ {
			buf[x] = byte(data[y][x])
		}
		file.Write(buf)
	}
}
