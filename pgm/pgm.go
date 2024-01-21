package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM represents a PGM image.
type PGM struct {
	data        [][]uint8
	width       int
	height      int
	magicNumber string
	max         int
}

// PBM represents a PBM image.
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// P2P5 represents the common fields for P2 and P5 formats.
type P2P5 struct {
	width, height int
	max           int
}

// PPM represents a colored PPM image.
type PPM struct {
	data          [][]Color // Assuming Color is a structure representing an RGB color.
	width, height int
	magicNumber   string
	max           int
}

// Color represents an RGB color.
type Color struct {
	R, G, B uint8
}

// NewPBM creates a new PBM image with the specified width and height.
func NewPBM(width, height int) *PBM {
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
	}
	return &PBM{
		data:   data,
		width:  width,
		height: height,
	}
}

// ReadPGM reads a PGM file and returns a PGM struct.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, errors.New("empty file")
	}

	pgm := &PGM{}
	pgm.magicNumber = scanner.Text()

	if !scanner.Scan() {
		return nil, errors.New("missing width and height")
	}
	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return nil, errors.New("invalid width and height")
	}

	pgm.width, err = strconv.Atoi(fields[0])
	if err != nil {
		return nil, errors.New("invalid width")
	}

	pgm.height, err = strconv.Atoi(fields[1])
	if err != nil {
		return nil, errors.New("invalid height")
	}

	if !scanner.Scan() {
		return nil, errors.New("missing maximum value")
	}
	line = scanner.Text()
	pgm.max, err = strconv.Atoi(line)
	if err != nil {
		return nil, errors.New("invalid maximum value")
	}

	pgm.data = make([][]uint8, pgm.height)
	for y := 0; y < pgm.height; y++ {
		if !scanner.Scan() {
			return nil, errors.New("missing image data")
		}
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) != pgm.width {
			return nil, errors.New("invalid image data")
		}

		pgm.data[y] = make([]uint8, pgm.width)
		for x, value := range fields {
			val, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return nil, errors.New("invalid pixel value")
			}
			pgm.data[y][x] = uint8(val)
		}
	}

	return pgm, nil
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

// Save saves the PGM image to a file.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write PGM header
	fmt.Fprintf(writer, "%s\n", pgm.magicNumber)
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	fmt.Fprintf(writer, "%d\n", pgm.max)

	// Write image data
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			fmt.Fprintf(writer, "%d ", pgm.data[y][x])
		}
		fmt.Fprintln(writer)
	}

	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pgm.data[y][x] = uint8(pgm.max - int(pgm.data[y][x]))
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		for x1, x2 := 0, pgm.width-1; x1 < x2; x1, x2 = x1+1, x2-1 {
			pgm.data[y][x1], pgm.data[y][x2] = pgm.data[y][x2], pgm.data[y][x1]
		}
	}
}

// Flop flips the PGM image vertically.
func (pgm *PGM) Flop() {
	for x := 0; x < pgm.width; x++ {
		for y1, y2 := 0, pgm.height-1; y1 < y2; y1, y2 = y1+1, y2-1 {
			pgm.data[y1][x], pgm.data[y2][x] = pgm.data[y2][x], pgm.data[y1][x]
		}
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the maximum value of the PGM image pixels.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	oldMax := pgm.max
	pgm.max = int(maxValue)

	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pgm.data[y][x] = uint8(float64(pgm.data[y][x]) * float64(pgm.max) / float64(oldMax))
		}
	}
}

// Rotate90CW rotates the PGM image 90 degrees clockwise.
func (pgm *PGM) Rotate90CW() {
	newWidth, newHeight := pgm.height, pgm.width
	newData := make([][]uint8, newHeight)

	for i := 0; i < newHeight; i++ {
		newData[i] = make([]uint8, newWidth)
	}

	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			newData[x][pgm.height-1-y] = pgm.data[y][x]
		}
	}

	pgm.width, pgm.height = newWidth, newHeight
	pgm.data = newData
}

// ToPBM converts the PGM image to a PBM image (black and white).
func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{}
	pbm.magicNumber = "P2"
	pbm.width = pgm.width
	pbm.height = pgm.height
	pbm.data = make([][]bool, pgm.height)

	for y := 0; y < pgm.height; y++ {
		pbm.data[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			pbm.data[y][x] = pgm.data[y][x] != 0
		}
	}

	return pbm
}

// NewPGM creates a new instance of the PGM structure with the specified dimensions.
func NewPGM(width, height, max int) *PGM {
	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}
	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: "P2",
		max:         max,
	}
}
