package pbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
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

	scanner := bufio.NewScanner(file)
	var magicNumber string

	// Read the magic number
	if scanner.Scan() {
		magicNumber = scanner.Text()
	} else {
		return nil, errors.New("missing magic number")
	}

	// Determine the width, height, and image format
	var width, height int
	var data [][]bool

	switch magicNumber {
	case "P1", "P4":
		width, height, data, err = readP1P4(scanner)
	case "P2", "P5":
		return nil, errors.New("P2 and P5 formats are not supported")
	default:
		return nil, errors.New("unsupported magic number")
	}

	if err != nil {
		return nil, err
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
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

	// Write the magic number
	fmt.Fprintln(writer, pbm.magicNumber)

	// Write the image data
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pixel {
				fmt.Fprint(writer, "1 ")
			} else {
				fmt.Fprint(writer, "0 ")
			}
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
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

// Helper function to reverse a boolean slice in place.
func reverseBoolSlice(slice []bool) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Helper function to read P1 and P4 formats.
func readP1P4(scanner *bufio.Scanner) (width, height int, data [][]bool, err error) {
	// Read the image size
	if !scanner.Scan() {
		return 0, 0, nil, errors.New("missing image size")
	}

	size := strings.Split(scanner.Text(), " ")
	if len(size) != 2 {
		return 0, 0, nil, errors.New("invalid image size format")
	}

	width, err = strconv.Atoi(size[0])
	if err != nil {
		return 0, 0, nil, errors.New("invalid width")
	}

	height, err = strconv.Atoi(size[1])
	if err != nil {
		return 0, 0, nil, errors.New("invalid height")
	}

	// Read the image data
	data = make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
		if !scanner.Scan() {
			return 0, 0, nil, errors.New("missing image data")
		}

		rowData := strings.Fields(scanner.Text())
		if len(rowData) != width {
			return 0, 0, nil, errors.New("invalid image data format")
		}

		for j, value := range rowData {
			data[i][j], err = strconv.ParseBool(value)
			if err != nil {
				return 0, 0, nil, errors.New("invalid pixel value")
			}
		}
	}

	return width, height, data, nil
}
