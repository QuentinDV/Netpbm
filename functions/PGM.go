package functions

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM struct to represent a PGM image
type PGM struct {
	data        [][]uint8
	width       int
	height      int
	magicNumber string
	max         int
}

// ReadPGM reads a PGM image from a file
func ReadPGM(fileName string) (*PGM, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	magicNumber := scanner.Text()
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, errors.New("uffgnsupported pgm format")
	}

	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(dimensions[0])
	height, _ := strconv.Atoi(dimensions[1])

	scanner.Scan()
	maxVal, _ := strconv.Atoi(scanner.Text())

	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
		for j := range data[i] {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())
			data[i][j] = uint8(val)
		}
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         maxVal,
	}, nil
}

// Size returns the width and height of the PGM image
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the pixel value at (x, y)
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the pixel value at (x, y)
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file
func (pgm *PGM) Save(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Fprintln(writer, pgm.data[i][j])
		}
	}

	return writer.Flush()
}

// Invert inverts the colors of the PGM image
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

// Flop flips the PGM image vertically
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i]
	}
}

// Rotate90CW rotates the PGM image 90 degrees clockwise
func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width)
	for i := range newData {
		newData[i] = make([]uint8, pgm.height)
		for j := range newData[i] {
			newData[i][j] = pgm.data[pgm.height-j-1][i]
		}
	}
	pgm.data = newData
	pgm.width, pgm.height = pgm.height, pgm.width
}

// SetMagicNumber sets the magic number of the PGM image
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// ToPBM converts the PGM image to PBM
func (pgm *PGM) ToPBM() *PBM {
	// This function requires implementation of PBM structure and conversion logic
	return &PBM{}
}

func PGMfunction(filename string) {
	// Example usage:
	pgm, err := ReadPGM("img/" + filename + ".pgm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(pgm.data)
	fmt.Printf("Magic Number: %s\n", pgm.magicNumber)
	fmt.Printf("Width: %d\n", pgm.width)
	pgm.Invert()
	fmt.Println(pgm.data)
	fmt.Printf("Height: %d\n", pgm.height)
	fmt.Printf("Max Value: %d\n", pgm.max)
	// fmt.Println(pgm.data[1][1])
}
