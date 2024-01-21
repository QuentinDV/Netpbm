package Netpbm

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// PPM represents a Portable PixMap image.
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint
}

// Pixel represents a color pixel with red, green, and blue components.
type Pixel struct {
	R, G, B uint8
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ppm := &PPM{}

	// Read and parse header
	scanner.Scan()
	ppm.magicNumber = scanner.Text()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &ppm.width, &ppm.height)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &ppm.max)

	// Read pixel data
	ppm.data = make([][]Pixel, ppm.height)
	for i := 0; i < ppm.height; i++ {
		ppm.data[i] = make([]Pixel, ppm.width)
		for j := 0; j < ppm.width; j++ {
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "%d %d %d", &ppm.data[i][j].R, &ppm.data[i][j].G, &ppm.data[i][j].B)
		}
	}

	return ppm, nil
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write header
	fmt.Fprintf(writer, "%s\n", ppm.magicNumber)
	fmt.Fprintf(writer, "%d %d\n", ppm.width, ppm.height)
	fmt.Fprintf(writer, "%d\n", ppm.max)

	// Write pixel data
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			fmt.Fprintf(writer, "%d %d %d\n", ppm.data[i][j].R, ppm.data[i][j].G, ppm.data[i][j].B)
		}
	}

	return writer.Flush()
}

// Invert inverts the colors of the PPM image.
// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width/2; j++ {
			ppm.data[i][j], ppm.data[i][ppm.width-j-1] = ppm.data[i][ppm.width-j-1], ppm.data[i][j]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	for i := 0; i < ppm.height/2; i++ {
		ppm.data[i], ppm.data[ppm.height-i-1] = ppm.data[ppm.height-i-1], ppm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = uint(maxValue)
}

// Rotate90CW rotates the PPM image 90Â° clockwise.
func (ppm *PPM) Rotate90CW() {
	newData := make([][]Pixel, ppm.width)
	for i := range newData {
		newData[i] = make([]Pixel, ppm.height)
	}

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			newData[j][ppm.height-1-i] = ppm.data[i][j]
		}
	}

	ppm.data = newData
	ppm.width, ppm.height = ppm.height, ppm.width
}

// ToPGM converts the PPM image to PGM.
// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	pgm := &PGM{
		data:        make([][]uint8, ppm.height),
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         ppm.max, // Convert ppm.max to uint
	}

	for i := 0; i < ppm.height; i++ {
		pgm.data[i] = make([]uint8, ppm.width)
		for j := 0; j < ppm.width; j++ {
			// Convert RGB to grayscale
			gray := uint8(0.299*float64(ppm.data[i][j].R) + 0.587*float64(ppm.data[i][j].G) + 0.114*float64(ppm.data[i][j].B))
			pgm.data[i][j] = gray
		}
	}

	return pgm
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		data:        make([][]bool, ppm.height),
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	for i := 0; i < ppm.height; i++ {
		pbm.data[i] = make([]bool, ppm.width)
		for j := 0; j < ppm.width; j++ {
			// Convert RGB to binary (black or white)
			gray := 0.299*float64(ppm.data[i][j].R) + 0.587*float64(ppm.data[i][j].G) + 0.114*float64(ppm.data[i][j].B)
			pbm.data[i][j] = gray > float64(ppm.max)/2
		}
	}

	return pbm
}

// Point represents a point in the image.
type Point struct {
	X, Y int
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))

	xIncrement := dx / float64(steps)
	yIncrement := dy / float64(steps)

	x, y := float64(p1.X), float64(p1.Y)

	for i := 0; i <= steps; i++ {
		ppm.Set(int(x), int(y), color)
		x += xIncrement
		y += yIncrement
	}
}

// DrawRectangle draws a rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X, p1.Y + height}
	p4 := Point{p1.X + width, p1.Y + height}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p4, color)
	ppm.DrawLine(p4, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for i := p1.Y; i < p1.Y+height; i++ {
		for j := p1.X; j < p1.X+width; j++ {
			ppm.Set(j, i, color)
		}
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	for i := 0; i <= radius*2; i++ {
		for j := 0; j <= radius*2; j++ {
			x := i - radius
			y := j - radius

			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
			}
		}
	}
}

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for i := -radius; i <= radius; i++ {
		for j := -radius; j <= radius; j++ {
			if i*i+j*j <= radius*radius {
				ppm.Set(center.X+i, center.Y+j, color)
			}
		}
	}
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle draws a filled triangle.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Sort points by Y coordinate
	points := []Point{p1, p2, p3}
	for i := 0; i < 2; i++ {
		for j := i + 1; j < 3; j++ {
			if points[i].Y > points[j].Y {
				points[i], points[j] = points[j], points[i]
			}
		}
	}

	// Calculate slopes
	slope1 := float64(points[1].X-points[0].X) / float64(points[1].Y-points[0].Y)
	slope2 := float64(points[2].X-points[0].X) / float64(points[2].Y-points[0].Y)
	slope3 := float64(points[2].X-points[1].X) / float64(points[2].Y-points[1].Y)

	// Draw upper part of the triangle
	for y := points[0].Y; y <= points[1].Y; y++ {
		x1 := int(float64(points[0].X) + slope1*float64(y-points[0].Y))
		x2 := int(float64(points[0].X) + slope2*float64(y-points[0].Y))
		ppm.DrawLine(Point{x1, y}, Point{x2, y}, color)
	}

	// Draw lower part of the triangle
	for y := points[1].Y + 1; y <= points[2].Y; y++ {
		x1 := int(float64(points[1].X) + slope3*float64(y-points[1].Y))
		x2 := int(float64(points[0].X) + slope2*float64(y-points[0].Y))
		ppm.DrawLine(Point{x1, y}, Point{x2, y}, color)
	}
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon draws a filled polygon.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Find bounding box
	minX, minY, maxX, maxY := points[0].X, points[0].Y, points[0].X, points[0].Y
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Fill the polygon using scanlines
	for y := minY; y <= maxY; y++ {
		intersections := 0
		for i := 0; i < len(points); i++ {
			x1, y1 := points[i].X, points[i].Y
			x2, y2 := points[(i+1)%len(points)].X, points[(i+1)%len(points)].Y

			// Check if the scanline intersects the edge
			if (y1 <= y && y < y2) || (y2 <= y && y < y1) {
				// Check if the scanline intersects the edge
				if (y1 <= y && y < y2) || (y2 <= y && y < y1) {
					// Calculate the intersection point
					xi := int(float64(x1) + (float64(y-y1)/float64(y2-y1))*float64(x2-x1))

					// If the intersection point is to the left of the current X position, increment the intersection count
					if xi < points[i].X {
						intersections++
					}
				}
			}

			// If the number of intersections is odd, color the pixel
			if intersections%2 != 0 {
				ppm.Set(points[i].X, y, color)
			}
		}
	}
}

// ToImage converts the PPM image to the Go image.Image interface.
func (ppm *PPM) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, ppm.width, ppm.height))

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			img.Set(x, y, color.RGBA{pixel.R, pixel.G, pixel.B, 255})
		}
	}

	return img
}

// SavePNG saves the PPM image as a PNG file.
func (ppm *PPM) SavePNG(filename string) error {
	img := ppm.ToImage()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
