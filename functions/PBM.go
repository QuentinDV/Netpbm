package functions

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func readP1(scanner *bufio.Scanner, width, height int) ([][]bool, error) {
	var data [][]bool
	for i := 0; i < height && scanner.Scan(); i++ {
		line := scanner.Text()
		row := make([]bool, width)

		for j, char := range line {
			if char == '1' {
				row[j] = true
			} else {
				row[j] = false
			}
		}

		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func readP4(file *os.File, width, height int) ([][]bool, error) {
	// Calculate the number of bytes needed for the data
	dataSize := (width + 7) / 8 * height
	data := make([]byte, dataSize)

	_, err := file.Read(data)
	if err != nil {
		return nil, err
	}

	var result [][]bool
	for i := 0; i < height; i++ {
		row := make([]bool, width)
		for j := 0; j < width; j++ {
			bit := (data[i*(width/8)+j/8] >> (7 - j%8)) & 1
			row[j] = bit == 1
		}
		result = append(result, row)
	}

	return result, nil
}

func (pbm *PBM) PrintData() {
	for i := 0; i < pbm.height; i++ {
		fmt.Println(pbm.data[i])
	}
}

func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the magic number, width, and height to the file
	fmt.Fprintf(writer, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Write P1 or P4 data based on the magic number
	if pbm.magicNumber == "P1" {
		err = writeP1Data(writer, pbm.data)
	} else if pbm.magicNumber == "P4" {
		err = writeP4Data(file, pbm.data)
	} else {
		return fmt.Errorf("unsupported PBM format: %s", pbm.magicNumber)
	}

	if err != nil {
		return err
	}

	return writer.Flush()
}

// writeP1Data writes P1 (ASCII) formatted data to the writer
func writeP1Data(writer *bufio.Writer, data [][]bool) error {
	var err error
	for _, row := range data {
		for _, pixel := range row {
			if pixel {
				_, err = writer.WriteString("1")
			} else {
				_, err = writer.WriteString("0")
			}
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// writeP4Data writes P4 (binary) formatted data to the file
func writeP4Data(file *os.File, data [][]bool) error {
	for _, row := range data {
		for _, pixel := range row {
			if pixel {
				_, err := file.Write([]byte{1})
				if err != nil {
					return err
				}
			} else {
				_, err := file.Write([]byte{0})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			if pbm.data[i][j] {
				pbm.data[i][j] = false
			} else {
				pbm.data[i][j] = true
			}
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		// Swap values from left to right
		for left, right := 0, pbm.width-1; left < right; left, right = left+1, right-1 {
			pbm.data[i][left], pbm.data[i][right] = pbm.data[i][right], pbm.data[i][left]
		}
	}
}

func (pbm *PBM) Flop() {
	// Swap rows from top to bottom
	for top, bottom := 0, pbm.height-1; top < bottom; top, bottom = top+1, bottom-1 {
		pbm.data[top], pbm.data[bottom] = pbm.data[bottom], pbm.data[top]
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func PBMfunction(filename string) {
	filePath := "img/" + filename + ".pbm"

	pbm, err := ReadPBM(filePath)
	if err != nil {
		fmt.Println("Error reading PBM image:", err)
		return
	}

	// Use the PBM struct as needed
	// pbm.PrintData()
	fmt.Printf("Magic Number: %s\n", pbm.magicNumber)
	fmt.Printf("Height: %d\n", pbm.height)
	fmt.Println(pbm.Size())
	fmt.Println(pbm.At(5, 6))
	pbm.Set(5, 6, false)
	fmt.Println(pbm.At(5, 6))
	pbm.Save("img/Save.pbm")

	// Inverted
	pbm.Invert()
	// pbm.PrintData()
	fmt.Printf("Magic Number: %s\n", pbm.magicNumber)
	pbm.SetMagicNumber("P4")
	fmt.Printf("Magic Number: %s\n", pbm.magicNumber)
	fmt.Printf("Height: %d\n", pbm.height)
	fmt.Println(pbm.Size())
	fmt.Println(pbm.At(5, 6))
	pbm.Set(5, 6, false)
	fmt.Println(pbm.At(5, 6))
	pbm.Save("img/SaveInvert.pbm")

}
