package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// PBM represents a PBM image.
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from the specified file.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read magic number
	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

	// Read dimensions
	dimensions, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading dimensions: %v", err)
	}
	var width, height int
	_, err = fmt.Sscanf(strings.TrimSpace(dimensions), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}

	data := make([][]bool, height)

	for i := range data {
		data[i] = make([]bool, width)
	}

	if magicNumber == "P1" {
		// Read P1 format (ASCII)
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index out of range at row %d", y)
				}
				data[y][x] = field == "1"
			}
		}

	} else if magicNumber == "P4" {
		// Read P4 format (binary)

		// Calculate the expected number of bytes per row for P4 format
		expectedBytesPerRow := (width + 7) / 8

		// Iterate over each row in the image
		for y := 0; y < height; y++ {
			// Read a row of bytes from the input file.
			row := make([]byte, expectedBytesPerRow)
			n, err := reader.Read(row)
			if err != nil {
				// Handle the case where unexpected end of file occurs
				if err == io.EOF {
					return nil, fmt.Errorf("unexpected end of file at row %d", y)
				}
				// Handle other errors while reading pixel data.
				return nil, fmt.Errorf("error reading pixel data at row %d: %v", y, err)
			}

			// Check if the actual number of bytes read matches the expected number.
			if n < expectedBytesPerRow {
				return nil, fmt.Errorf("unexpected end of file at row %d, expected %d bytes, got %d", y, expectedBytesPerRow, n)
			}

			// Iterate over each pixel in the row
			for x := 0; x < width; x++ {
				// Calculate the index of the byte containing the current pixel
				byteIndex := x / 8

				// Calculate the bit index within the byte for the current pixel
				bitIndex := 7 - (x % 8)

				// Extract the bit value from the byte using a bitwise AND operation
				// Shift the decimal value to the right by the bit index, then perform a bitwise AND with 1
				bitValue := (int(row[byteIndex]) >> bitIndex) & 1

				// Set the corresponding pixel in the image data based on the extracted bit value
				data[y][x] = bitValue != 0
			}
		}
	}

	return &PBM{data, width, height, magicNumber}, nil
}

// Size returns the width and height of the image
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y)
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set sets the value of the pixel at (x, y)
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save saves the PBM image to the specified file.
func (pbm *PBM) Save(filename string) error {
	if pbm == nil {
		return errors.New("cannot save a nil PBM")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write magic number, width, and height
	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Choose the appropriate method based on the magic number
	if pbm.magicNumber == "P1" {
		return pbm.saveP1(file)
	} else if pbm.magicNumber == "P4" {
		return pbm.saveP4(file)
	} else {
		return fmt.Errorf("unsupported magic number: %s", pbm.magicNumber)
	}
}

// saveP1 saves the PBM image in P1 format (ASCII)
func (pbm *PBM) saveP1(file *os.File) error {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			// Write the binary value of the pixel
			if pbm.data[i][j] {
				fmt.Fprint(file, "1")
			} else {
				fmt.Fprint(file, "0")
			}

			// Add a space after each pixel, except the last one in a row
			if j < pbm.width-1 {
				fmt.Fprint(file, " ")
			}
		}
		// Add a newline after each row
		fmt.Fprintln(file)
	}
	return nil
}

// saveP4 saves the PBM image in P4 format (binary)
func (pbm *PBM) saveP4(file *os.File) error {
	expectedBytesPerRow := (pbm.width + 7) / 8
	for y := 0; y < pbm.height; y++ {
		row := make([]byte, expectedBytesPerRow)
		for x := 0; x < pbm.width; x++ {
			byteIndex := x / 8
			bitIndex := 7 - (x % 8)
			if pbm.data[y][x] {
				row[byteIndex] |= 1 << bitIndex
			}
		}
		_, err := file.Write(row)
		if err != nil {
			return fmt.Errorf("error writing pixel data at row %d: %v", y, err)
		}
	}
	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width/2; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
