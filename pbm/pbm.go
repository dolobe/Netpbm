package pbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data        [][]bool
	width       int
	height      int
	magicNumber string
}

// ReadPBM lit une image PBM à partir d'un fichier et renvoie une struct qui représente l'image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	if !scanner.Scan() {
		return nil, errors.New("missing magic number")
	}
	magicNumber := scanner.Text()

	// Read width and height
	if !scanner.Scan() {
		return nil, errors.New("missing width")
	}
	width, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, err
	}

	if !scanner.Scan() {
		return nil, errors.New("missing height")
	}
	height, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, err
	}

	// Read pixel values
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
		lineValues := strings.Fields(scanner.Text())
		if len(lineValues) != width {
			fmt.Printf("Error: line %d has %d values instead of %d\n", i+1, len(lineValues), width)
			return nil, errors.New("incorrect number of pixel values in the line")
		}
		for j := range data[i] {
			value, err := strconv.Atoi(lineValues[j])
			if err != nil {
				return nil, err
			}
			data[i][j] = value != 0
		}
		if !scanner.Scan() {
			fmt.Println("Error: pixel data missing at line", i+2)
			return nil, errors.New("missing pixel data")
		}
	}

	// Create and return the PBM structure
	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At renvoie la valeur du pixel à la position (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set définit la valeur du pixel à la position (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save enregistre l'image PBM dans un fichier et renvoie une erreur en cas de problème.
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

// Invert inverse les couleurs de l'image PBM.
func (pbm *PBM) Invert() {
	for i := range pbm.data {
		for j := range pbm.data[i] {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip retourne l'image PBM horizontalement.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width/2; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}

// Flop retourne l'image PBM verticalement.
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j], pbm.data[pbm.height-i-1][j] = pbm.data[pbm.height-i-1][j], pbm.data[i][j]
		}
	}
}

// SetMagicNumber définit le nombre magique de l'image PBM.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func main() {
	// Example usage:
	pbm, err := ReadPBM("testP1.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	width, height := pbm.Size()
	fmt.Printf("Image size: %d x %d\n", width, height)

	// Enregistrer l'image modifier
	err = pbm.Save("modified.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
