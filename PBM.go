package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PBM represents a PBM image.
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	// Ouvrir le fichier
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	MagicNumber := scanner.Text()

	// Ignore comments
	for strings.HasPrefix(scanner.Text(), "#") {
		scanner.Scan()
	}

	// Read width and height
	scanner.Scan()
	var Width, Height int
	fmt.Sscanf(scanner.Text(), "%d %d", &Width, &Height)

	// Read P1 or P4 data based on the magic number
	var Data [][]bool
	if MagicNumber == "P1" {
		Data, err = readP1(scanner, Width, Height)
	} else if MagicNumber == "P4" {
		Data, err = readP4(file, Width, Height)
	} else {
		return nil, fmt.Errorf("unsupported PBM format: %s", MagicNumber)
	}

	if err != nil {
		return nil, err
	}

	return &PBM{
		magicNumber: MagicNumber,
		width:       Width,
		height:      Height,
		data:        Data,
	}, nil
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write magic number, width, and height
	fmt.Fprintf(writer, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Write image data
	if pbm.magicNumber == "P1" {
		writeP1(writer, pbm.data)
	} else if pbm.magicNumber == "P4" {
		writeP4(file, pbm.data)
	} else {
		return fmt.Errorf("unsupported PBM format: %s", pbm.magicNumber)
	}

	writer.Flush()
	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width/2; x++ {
			pbm.data[y][x], pbm.data[y][pbm.width-x-1] = pbm.data[y][pbm.width-x-1], pbm.data[y][x]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	for y := 0; y < pbm.height/2; y++ {
		pbm.data[y], pbm.data[pbm.height-y-1] = pbm.data[pbm.height-y-1], pbm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

// Additional helper functions

func readP1(scanner *bufio.Scanner, width, height int) ([][]bool, error) {
	var data [][]bool

	for y := 0; y < height; y++ {
		scanner.Scan()
		line := scanner.Text()
		var row []bool
		for _, char := range line {
			if char == '0' {
				row = append(row, false)
			} else if char == '1' {
				row = append(row, true)
			}
		}
		data = append(data, row)
	}

	return data, nil
}

func readP4(file *os.File, width, height int) ([][]bool, error) {
	data := make([][]bool, height)
	for y := 0; y < height; y++ {
		data[y] = make([]bool, width)
		buf := make([]byte, width/8)
		_, err := file.Read(buf)
		if err != nil {
			return nil, err
		}
		for x := 0; x < width; x++ {
			// Extract the bit from the byte
			bit := (buf[x/8] >> uint(7-x%8)) & 1
			data[y][x] = bit != 0
		}
	}
	return data, nil
}

func writeP1(writer *bufio.Writer, data [][]bool) {
	for y := 0; y < len(data); y++ {
		for x := 0; x < len(data[y]); x++ {
			value := 0
			if data[y][x] {
				value = 1
			}
			fmt.Fprintf(writer, "%d ", value)
		}
		fmt.Fprintln(writer)
	}
}

func writeP4(file *os.File, data [][]bool) {
	for y := 0; y < len(data); y++ {
		var byteValue byte
		for x := 0; x < len(data[y]); x++ {
			// Set the corresponding bit in the byte
			bit := byte(0)
			if data[y][x] {
				bit = 1
			}
			byteValue = (byteValue << 1) | bit
		}
		file.Write([]byte{byteValue})
	}
}
