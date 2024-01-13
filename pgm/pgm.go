package pgm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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

	// Read width and height
	scanner.Scan()
	width, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	height, _ := strconv.Atoi(scanner.Text())

	// Read max value
	scanner.Scan()
	maxValue, _ := strconv.Atoi(scanner.Text())

	// Read pixel values
	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
		for j := range data[i] {
			scanner.Scan()
			value, _ := strconv.Atoi(scanner.Text())
			data[i][j] = uint8(value)
		}
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         maxValue,
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

	// Write pixel values
	for _, row := range pgm.data {
		for _, value := range row {
			fmt.Fprintf(writer, "%d\n", value)
		}
	}

	writer.Flush()

	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := range pgm.data {
		for j := range pgm.data[i] {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for i := range pgm.data {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i]
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
	// Create a new PGM image with swapped width and height
	newData := make([][]uint8, pgm.width)
	for i := range newData {
		newData[i] = make([]uint8, pgm.height)
	}

	// Rotate the pixel values
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			newData[j][pgm.height-i-1] = pgm.data[i][j]
		}
	}

	// Update the PGM image with rotated data
	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = newData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	// Define a threshold value (you can adjust this based on your needs)
	threshold := uint8(128)

	// Create a new PBM image
	pbmData := make([][]bool, pgm.height)
	for i := range pbmData {
		pbmData[i] = make([]bool, pgm.width)
	}

	// Convert pixel values based on the threshold
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pbmData[i][j] = pgm.data[i][j] >= threshold
		}
	}

	// Create and return the new PBM image
	return &PBM{
		data:        pbmData,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1", // Assuming P1 format for binary PBM
	}
}

// PBM structure for ToPBM conversion
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
	pgm           *PGM // Embed PGM structure
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	scanner.Scan()
	magicNumber := scanner.Text()

	// Read width and height
	scanner.Scan()
	width, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	height, _ := strconv.Atoi(scanner.Text())

	// Read pixel values
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
		for j := range data[i] {
			scanner.Scan()
			value, _ := strconv.Atoi(scanner.Text())
			data[i][j] = value != 0
		}
	}

	// Create and return the PBM structure with embedded PGM structure
	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		pgm:         nil, // Initialize as needed
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

	// Write pixel values
	for _, row := range pbm.data {
		for _, value := range row {
			if value {
				fmt.Fprint(writer, "1\n")
			} else {
				fmt.Fprint(writer, "0\n")
			}
		}
	}

	writer.Flush()

	return nil
}
func main() {
	pgm, err := ReadPGM(".pgm")
	if err != nil {
		fmt.Println("Error reading PGM file:", err)
		return
	}

	pgm.Invert()
	pgm.Save("inverted.pgm")

	pgm.Flip()
	pgm.Save("flipped.pgm")

	pgm.Flop()
	pgm.Save("flopped.pgm")
}
