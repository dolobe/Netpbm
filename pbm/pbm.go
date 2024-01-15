package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var magicNumber string
	var width, height int
	var data [][]bool

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue // Skip comments
		}

		fields := strings.Fields(line)

		if magicNumber == "" {
			magicNumber = fields[0]
		} else if width == 0 {
			width, err = strconv.Atoi(fields[0])
			if err != nil {
				return nil, err
			}
		} else if height == 0 {
			height, err = strconv.Atoi(fields[0])
			if err != nil {
				return nil, err
			}
			break
		}
	}

	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("unsupported magic number: %s", magicNumber)
	}

	data = make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
	}

	var i, j int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue // Skip comments
		}

		fields := strings.Fields(line)
		for _, val := range fields {
			if j >= width {
				j = 0
				i++
			}
			if i >= height {
				break
			}

			b, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			data[i][j] = b != 0
			j++
		}
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

	// Write the PBM header
	_, err = file.WriteString(pbm.magicNumber + "\n")
	if err != nil {
		return err
	}

	// Write the dimensions of the image
	_, err = fmt.Fprintf(file, "%d %d\n", pbm.width, pbm.height)
	if err != nil {
		return err
	}

	// Write the image data
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			if pbm.At(x, y) {
				_, err = file.WriteString("1 ")
			} else {
				_, err = file.WriteString("0 ")
			}
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintln(file) // newline after each row
		if err != nil {
			return err
		}
	}

	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for i := range pbm.data {
		for j := range pbm.data[i] {
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
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j], pbm.data[pbm.height-i-1][j] = pbm.data[pbm.height-i-1][j], pbm.data[i][j]
		}
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func main() {
	// Example usage:
	pbm, err := ReadPBM("../testImages/pbm/testP1.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	width, height := pbm.Size()
	fmt.Printf("Image size: %d x %d\n", width, height)

	// Do other operations...

	// Save the modified image
	err = pbm.Save("modified.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
