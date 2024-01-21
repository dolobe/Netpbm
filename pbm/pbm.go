package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PBM represents a Portable BitMap image.
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// NewPBM creates a new PBM image with the specified width and height.
func NewPBM(width, height int) *PBM {
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
	}
	return &PBM{
		magicNumber: "P1",
		data:        data,
		width:       width,
		height:      height,
	}
}

// ReadPBM reads a PBM image from a file and returns a struct representing the image.
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
		bytesPerRow := (width + 7) / 8 // Number of bytes needed to store a row
		for y := 0; y < height; y++ {
			row := make([]byte, bytesPerRow)
			_, err := reader.Read(row)
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			for x := 0; x < width; x++ {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				data[y][x] = (row[byteIndex]>>bitIndex)&1 == 1
			}
		}
	}

	// Create and return the PBM struct
	pbmImage := &PBM{
		magicNumber: magicNumber,
		width:       width,
		height:      height,
		data:        data,
	}

	return pbmImage, nil
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

// Save saves the PBM image to a file.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write magic number
	_, err = fmt.Fprintf(writer, "%s\n", pbm.magicNumber)
	if err != nil {
		return fmt.Errorf("error writing magic number: %v", err)
	}

	// Write dimensions
	_, err = fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height)
	if err != nil {
		return fmt.Errorf("error writing dimensions: %v", err)
	}

	if pbm.magicNumber == "P1" {
		// Write P1 format (ASCII)
		for y := 0; y < pbm.height; y++ {
			for x := 0; x < pbm.width; x++ {
				val := 0
				if pbm.data[y][x] {
					val = 1
				}
				_, err := fmt.Fprintf(writer, "%d ", val)
				if err != nil {
					return fmt.Errorf("error writing data at row %d, column %d: %v", y, x, err)
				}
			}
			_, err := fmt.Fprint(writer, "\n")
			if err != nil {
				return fmt.Errorf("error writing newline at row %d: %v", y, err)
			}
		}
	} else if pbm.magicNumber == "P4" {
		// Write P4 format (binary)
		bytesPerRow := (pbm.width + 7) / 8
		for y := 0; y < pbm.height; y++ {
			row := make([]byte, bytesPerRow)
			for x := 0; x < pbm.width; x++ {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				if pbm.data[y][x] {
					row[byteIndex] |= 1 << bitIndex
				}
			}
			_, err := writer.Write(row)
			if err != nil {
				return fmt.Errorf("error writing data at row %d: %v", y, err)
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing writer: %v", err)
	}

	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for y, row := range pbm.data {
		for x := range row {
			pbm.data[y][x] = !pbm.data[y][x]
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
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j], pbm.data[pbm.height-i-1][j] = pbm.data[pbm.height-i-1][j], pbm.data[i][j]
		}
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
