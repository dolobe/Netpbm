package Netpbm

import (
	"testing"
)

const imagePGMWidth = 15
const imagePGMHeight = 15
const imagePGMMax = 11

var testData = []int{
	11, 11, 11, 11, 11, 11, 11, 0, 0, 0, 0, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 8, 11, 0, 0, 0, 11,
	11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 5, 5, 0, 11, 11, 11, 11, 11, 0, 0, 11, 11, 11, 11, 11, 5, 0, 0, 0, 0, 11, 11, 11, 11, 0, 0, 11, 11, 11, 0, 0, 0, 11, 0, 7, 0, 0, 11, 11, 11, 0, 11, 11, 11, 0, 11, 11, 11, 0, 7, 11, 11, 0,
	0, 0, 11, 11, 11, 11, 0, 11, 11, 11, 0, 7, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11, 0, 7, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 0, 0, 11, 11, 11, 11,
	11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 0, 0, 7, 7, 7, 7, 7, 0, 0, 11, 11, 11, 11, 11, 11, 11, 11, 0, 0, 0, 0, 0, 0, 11, 11, 11, 11, 11,
}

var testInvertPGM = []int{
	0, 0, 0, 0, 0, 0, 0, 11, 11, 11, 11, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 11, 0, 0, 0,
	0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 11, 0, 0,
	0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 3, 0, 11, 11, 11,
	0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 6, 6, 11,
	0, 0, 0, 0, 0, 11, 11, 0, 0, 0, 0, 0, 6, 11, 11,
	11, 11, 0, 0, 0, 0, 11, 11, 0, 0, 0, 11, 11, 11, 0,
	11, 4, 11, 11, 0, 0, 0, 11, 0, 0, 0, 11, 0, 0, 0,
	11, 4, 0, 0, 11, 11, 11, 0, 0, 0, 0, 11, 0, 0, 0,
	11, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 0, 0,
	11, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 0, 0,
	0, 11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 0, 0,
	0, 11, 11, 0, 0, 0, 0, 0, 0, 0, 0, 11, 0, 0, 0,
	0, 0, 11, 11, 4, 4, 4, 4, 4, 11, 11, 0, 0, 0, 0,
	0, 0, 0, 0, 11, 11, 11, 11, 11, 11, 0, 0, 0, 0, 0,
}

var testFlipPGM = []int{
	11, 11, 11, 11, 0, 0, 0, 0, 11, 11, 11, 11, 11, 11, 11,
	11, 11, 11, 0, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11,
	11, 11, 0, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11,
	0, 0, 0, 11, 8, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11,
	0, 5, 5, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11,
	0, 0, 5, 11, 11, 11, 11, 11, 0, 0, 11, 11, 11, 11, 11,
	11, 0, 0, 0, 11, 11, 11, 0, 0, 11, 11, 11, 11, 0, 0,
	11, 11, 11, 0, 11, 11, 11, 0, 11, 11, 11, 0, 0, 7, 0,
	11, 11, 11, 0, 11, 11, 11, 11, 0, 0, 0, 11, 11, 7, 0,
	11, 11, 0, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 7, 0,
	11, 11, 0, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 7, 0,
	11, 11, 0, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11,
	11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 11, 11, 0, 0, 11,
	11, 11, 11, 11, 0, 0, 7, 7, 7, 7, 7, 0, 0, 11, 11,
	11, 11, 11, 11, 11, 0, 0, 0, 0, 0, 0, 11, 11, 11, 11,
}

var testFlopPGM = []int{
	11, 11, 11, 11, 0, 0, 0, 0, 0, 0, 11, 11, 11, 11, 11,
	11, 11, 0, 0, 7, 7, 7, 7, 7, 0, 0, 11, 11, 11, 11,
	11, 0, 0, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11, 11,
	11, 0, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11,
	0, 7, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11,
	0, 7, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 11, 11,
	0, 7, 11, 11, 0, 0, 0, 11, 11, 11, 11, 0, 11, 11, 11,
	0, 7, 0, 0, 11, 11, 11, 0, 11, 11, 11, 0, 11, 11, 11,
	0, 0, 11, 11, 11, 11, 0, 0, 11, 11, 11, 0, 0, 0, 11,
	11, 11, 11, 11, 11, 0, 0, 11, 11, 11, 11, 11, 5, 0, 0,
	11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 5, 5, 0,
	11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 8, 11, 0, 0, 0,
	11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 0, 11, 11,
	11, 11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 0, 11, 11, 11,
	11, 11, 11, 11, 11, 11, 11, 0, 0, 0, 0, 11, 11, 11, 11,
}

var testRotate90PGM = []int{
	11, 11, 11, 11, 0, 0, 0, 0, 0, 11, 11, 11, 11, 11, 11,
	11, 11, 0, 0, 7, 7, 7, 7, 0, 11, 11, 11, 11, 11, 11,
	11, 0, 0, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 11,
	11, 0, 11, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 11,
	0, 7, 11, 11, 11, 11, 0, 11, 11, 11, 11, 11, 11, 11, 11,
	0, 7, 11, 11, 11, 11, 0, 11, 11, 0, 0, 0, 0, 11, 11,
	0, 7, 11, 11, 11, 11, 0, 11, 0, 0, 11, 11, 11, 0, 11,
	0, 7, 11, 11, 11, 11, 11, 0, 0, 11, 11, 11, 11, 11, 0,
	0, 7, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0,
	0, 0, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 0,
	11, 0, 11, 11, 11, 11, 11, 11, 11, 11, 11, 8, 11, 11, 0,
	11, 11, 0, 11, 11, 11, 0, 0, 0, 11, 11, 11, 11, 0, 11,
	11, 11, 11, 0, 0, 0, 11, 11, 0, 5, 5, 0, 0, 11, 11,
	11, 11, 11, 11, 11, 11, 11, 11, 0, 0, 5, 0, 11, 11, 11,
	11, 11, 11, 11, 11, 11, 11, 11, 11, 0, 0, 0, 11, 11, 11,
}

func TestReadPGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	if pgm.magicNumber != "P2" {
		t.Error("Magic number not read correctly")
	}
	if pgm.width != imagePGMWidth {
		t.Error("Width not read correctly")
	}
	if pgm.height != imagePGMHeight {
		t.Error("Height not read correctly")
	}
	if pgm.max != imagePGMMax {
		t.Error("Max value not read correctly")
	}
	for i := 0; i < imageWidth*imageHeight; i++ {
		x := i % imageWidth
		y := i / imageWidth
		if pgm.data[y][x] != uint8(testData[i]) {
			t.Errorf("Pixel at (%d, %d) not read correctly", x, y)
		}
	}
	pgm, err = ReadPGM("./testImages/pgm/testP5.pgm")
}

func TestSizePGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	w, h := pgm.Size()
	if w != imagePGMWidth {
		t.Error("Width not read correctly")
	}
	if h != imagePGMHeight {
		t.Error("Height not read correctly")
	}
}

func TestAtPGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	if pgm.At(0, 8) != 0 {
		t.Error("Wrong value")
	}
}

func TestSetPGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.Set(0, 8, 5)
	if pgm.At(0, 8) != 5 {
		t.Error("Wrong value")
	}
}

func TestSavePGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.SetMagicNumber("P2")
	pgm.Save("./testImages/pgm/testP2a.pgm")
	pgm, err = ReadPGM("./testImages/pgm/testP2a.pgm")
	if err != nil {
		t.Error(err)
	}
	if pgm.magicNumber != "P2" {
		t.Error("Magic number not read correctly")
	}
	if pgm.width != imagePGMWidth {
		t.Error("Width not read correctly")
	}
	if pgm.height != imagePGMHeight {
		t.Error("Height not read correctly")
	}
	if pgm.max != imagePGMMax {
		t.Error("Max value not read correctly")
	}
	for i := 0; i < imageWidth*imageHeight; i++ {
		x := i % imageWidth
		y := i / imageWidth
		if pgm.data[y][x] != uint8(testData[i]) {
			t.Errorf("Pixel at (%d, %d) not read correctly", x, y)
		}
	}
	pgm, err = ReadPGM("./testImages/pgm/testP5.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.SetMagicNumber("P5")
	pgm.Save("./testImages/pgm/testP5a.pgm")
	pgm, err = ReadPGM("./testImages/pgm/testP5a.pgm")
	if err != nil {
		t.Error(err)
	}
	if pgm.magicNumber != "P5" {
		t.Error("Magic number not read correctly")
	}
	if pgm.width != imagePGMWidth {
		t.Error("Width not read correctly")
	}
	if pgm.height != imagePGMHeight {
		t.Error("Height not read correctly")
	}
	if pgm.max != imagePGMMax {
		t.Error("Max value not read correctly")
	}
	for i := 0; i < imageWidth*imageHeight; i++ {
		x := i % imageWidth
		y := i / imageWidth
		if pgm.data[y][x] != uint8(testData[i]) {
			t.Errorf("Pixel at (%d, %d) not read correctly", x, y)
		}
	}
}

func TestInvertPGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.Invert()
	for i := 0; i < imageWidth*imageHeight; i++ {
		x := i % imageWidth
		y := i / imageWidth
		if pgm.data[y][x] != uint8(testInvertPGM[i]) {
			t.Errorf("Pixel at (%d, %d) not read correctly", x, y)
		}
	}
}

func TestFlipPGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.Flip()
	for i := 0; i < imageWidth*imageHeight; i++ {
		x := i % imageWidth
		y := i / imageWidth
		if pgm.data[y][x] != uint8(testFlipPGM[i]) {
			t.Errorf("Pixel at (%d, %d) not read correctly", x, y)
		}
	}
}

func TestFlopPGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.Flop()
	for i := 0; i < imageWidth*imageHeight; i++ {
		x := i % imageWidth
		y := i / imageWidth
		if pgm.data[y][x] != uint8(testFlopPGM[i]) {
			t.Errorf("Pixel at (%d, %d) not read correctly", x, y)
		}
	}
}

func TestRotate90CWPGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.Rotate90CW()
	for i := 0; i < imageWidth*imageHeight; i++ {
		x := i % imageWidth
		y := i / imageWidth
		if pgm.data[y][x] != uint8(testRotate90PGM[i]) {
			t.Errorf("Pixel at (%d, %d) not read correctly", x, y)
		}
	}
}

func TestSetMagicNumberPGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.SetMagicNumber("P5")
	if pgm.magicNumber != "P5" {
		t.Error("Magic number not set correctly")
	}
}

func TestSetMaxValuePGM(t *testing.T) {
	pgm, err := ReadPGM("./testImages/pgm/testP2.pgm")
	if err != nil {
		t.Error(err)
	}
	pgm.SetMaxValue(5)
	if pgm.max != 5 {
		t.Error("Max value not set correctly")
	}
	for i := 0; i < imageWidth*imageHeight; i++ {
		x := i % imageWidth
		y := i / imageWidth
		if pgm.data[y][x] != uint8(float64(testData[i])/float64(5)*255) {
			t.Errorf("Pixel at (%d, %d) not read correctly, expected %d, got %d", x, y, uint8(float64(testData[i])/float64(5)*255), pgm.data[y][x])
		}
	}
}
