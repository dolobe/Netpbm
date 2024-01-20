package ppm

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"
	"strconv"
)

// PPM represents a PPM image.
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

// Pixel represents a color pixel.
type Pixel struct {
	R, G, B uint8
}

// Point represents a point in the image.
type Point struct {
	X, Y int
}

// PGM represents a Portable GrayMap image.
type PGM struct {
	magicNumber string
	data        [][]uint8
	width       int
	height      int
	max         uint8
}

// PBM represents a Portable BitMap image.
type PBM struct {
	magicNumber string
	data        [][]bool
	width       int
	height      int
}

// NewPGM creates a new PGM image with the specified width and height.
func NewPGM(width, height int) *PGM {
	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}
	return &PGM{
		data:   data,
		width:  width,
		height: height,
		max:    255, // Default max value for PGM
	}
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

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	ppm := &PPM{}

	// Read magic number
	scanner.Scan()
	ppm.magicNumber = scanner.Text()

	// Read width and height
	scanner.Scan()
	ppm.width, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	ppm.height, _ = strconv.Atoi(scanner.Text())

	// Read max value
	scanner.Scan()
	ppm.max, _ = strconv.Atoi(scanner.Text())

	// Initialize data
	ppm.data = make([][]Pixel, ppm.height)
	for i := range ppm.data {
		ppm.data[i] = make([]Pixel, ppm.width)
	}

	// Read pixel data
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			scanner.Scan()
			r, _ := strconv.Atoi(scanner.Text())
			scanner.Scan()
			g, _ := strconv.Atoi(scanner.Text())
			scanner.Scan()
			b, _ := strconv.Atoi(scanner.Text())
			ppm.data[y][x] = Pixel{uint8(r), uint8(g), uint8(b)}
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

	// Write magic number, width, height, and max value
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	// Write pixel data
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			fmt.Fprintf(writer, "%d %d %d ", ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B)
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			ppm.data[y][x].R = uint8(ppm.max) - ppm.data[y][x].R
			ppm.data[y][x].G = uint8(ppm.max) - ppm.data[y][x].G
			ppm.data[y][x].B = uint8(ppm.max) - ppm.data[y][x].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x1, x2 := 0, ppm.width-1; x1 < x2; x1, x2 = x1+1, x2-1 {
			ppm.data[y][x1], ppm.data[y][x2] = ppm.data[y][x2], ppm.data[y][x1]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	for y1, y2 := 0, ppm.height-1; y1 < y2; y1, y2 = y1+1, y2-1 {
		ppm.data[y1], ppm.data[y2] = ppm.data[y2], ppm.data[y1]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = int(maxValue)
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW() {
	// Create a new PPM image with swapped width and height
	newPPM := NewPPM(ppm.height, ppm.width)

	// Copy data to the new image, rotating it
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			newPPM.data[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	// Update the original image
	ppm.width, ppm.height = newPPM.width, newPPM.height
	ppm.data = newPPM.data
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Create a new PGM image with the same size
	pgm := NewPGM(ppm.width, ppm.height)

	// Convert color pixels to grayscale and set them in the new image
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			gray := uint8((int(ppm.data[y][x].R) + int(ppm.data[y][x].G) + int(ppm.data[y][x].B)) / 3)
			pgm.data[y][x] = gray
		}
	}

	return pgm
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	// Create a new PBM image with the same size
	pbm := NewPBM(ppm.width, ppm.height)

	// Convert color pixels to binary and set them in the new image
	for y := 0; y > ppm.height; y++ {
		for x := 0; x > ppm.width; x++ {
			// Assume that a pixel is black if at least one color channel is non-zero
			black := ppm.data[y][x].R != 0 || ppm.data[y][x].G != 0 || ppm.data[y][x].B != 0
			pbm.data[y][x] = black
		}
	}

	return pbm
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	x, y := p1.X, p1.Y
	sx, sy := 1, 1

	if dx < 0 {
		sx = -1
		dx = -dx
	}
	if dy < 0 {
		sy = -1
		dy = -dy
	}

	err := dx - dy

	for {
		ppm.Set(x, y, color)

		if x == p2.X && y == p2.Y {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x += sx
		}

		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// DrawRectangle draws a rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	// Draw the four sides of the rectangle using DrawLine.
	p2 := Point{p1.X + width - 1, p1.Y}
	p3 := Point{p1.X + width - 1, p1.Y + height - 1}
	p4 := Point{p1.X, p1.Y + height - 1}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	// Fill the rectangle by setting each pixel inside the rectangle to the specified color.
	for y := p1.Y; y < p1.Y+height; y++ {
		for x := p1.X; x < p1.X+width; x++ {
			ppm.Set(x, y, color)
		}
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	x, y, err := radius, 0, 0

	for x >= y {
		ppm.Set(center.X+x, center.Y+y, color)
		ppm.Set(center.X+y, center.Y+x, color)
		ppm.Set(center.X-y, center.Y+x, color)
		ppm.Set(center.X-x, center.Y+y, color)
		ppm.Set(center.X-x, center.Y-y, color)
		ppm.Set(center.X-y, center.Y-x, color)
		ppm.Set(center.X+y, center.Y-x, color)
		ppm.Set(center.X+x, center.Y-y, color)

		if err <= 0 {
			y++
			err += 2*y + 1
		}

		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// Fill the circle by setting each pixel inside the circle to the specified color.
	x, y := -radius, 0
	err, delta := 2-2*radius, 0

	for x < 0 {
		if center.X-x >= 0 && center.X-x < ppm.width && center.Y+y >= 0 && center.Y+y < ppm.height {
			ppm.Set(center.X-x, center.Y+y, color)
		}
		if center.X-x >= 0 && center.X-x < ppm.width && center.Y-y >= 0 && center.Y-y < ppm.height {
			ppm.Set(center.X-x, center.Y-y, color)
		}
		if center.X+x >= 0 && center.X+x < ppm.width && center.Y-y >= 0 && center.Y-y < ppm.height {
			ppm.Set(center.X+x, center.Y-y, color)
		}
		if center.X+x >= 0 && center.X+x < ppm.width && center.Y+y >= 0 && center.Y+y < ppm.height {
			ppm.Set(center.X+x, center.Y+y, color)
		}

		delta = 2*(err+y) - 1
		if err < 0 && delta <= 0 {
			x++
			err += x*2 + 1
			continue
		}

		delta = 2*(err-x) - 1
		if err > 0 && delta > 0 {
			y--
			err += 1 - y*2
			continue
		}

		x++
		err += x*2 + 1
		y--
		err += 1 - y*2
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
	// Fill the triangle by drawing three lines to create three smaller triangles.
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// Draw lines connecting consecutive points to form the polygon.
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	// Connect the last and first points to close the polygon.
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon draws a filled polygon.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	minY, maxY := ppm.height, 0

	// Find the bounding box of the polygon.
	for _, p := range points {
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Create a list to store the intersection points with each scanline.
	intersections := make([][]int, ppm.height)

	for i := range intersections {
		intersections[i] = make([]int, 0)
	}

	// Iterate through each edge of the polygon and find intersections with scanlines.
	for i := 0; i < len(points); i++ {
		p1, p2 := points[i], points[(i+1)%len(points)]
		ppm.findIntersections(p1, p2, &intersections)
	}

	// Fill the polygon row by row.
	for y := minY; y <= maxY; y++ {
		// Sort the intersection points based on the X-coordinate.
		sort.Ints(intersections[y])

		// Draw lines connecting consecutive intersection points.
		for i := 0; i < len(intersections[y])-1; i += 2 {
			ppm.DrawLine(Point{intersections[y][i], y}, Point{intersections[y][i+1], y}, color)
		}
	}
}

// findIntersections finds intersections between the polygon edges and a horizontal scanline.
func (ppm *PPM) findIntersections(p1, p2 Point, intersections *[][]int) {
	// Check if the edge intersects with the scanline.
	if p1.Y == p2.Y {
		return
	}
	if p1.Y > p2.Y {
		p1, p2 = p2, p1
	}

	x1, y1, x2, y2 := p1.X, p1.Y, p2.X, p2.Y

	if y1 >= ppm.height || y2 < 0 {
		return
	}

	if y1 < 0 {
		// Clip the edge to the upper edge of the image.
		x1 = x1 + (0-y1)*(x2-x1)/(y2-y1)
		y1 = 0
	}

	if y2 >= ppm.height {
		// Clip the edge to the lower edge of the image.
		x2 = x2 - (y2-ppm.height+1)*(x2-x1)/(y2-y1)
		y2 = ppm.height - 1
	}

	// Add the intersection points to the list.
	m := (x2 - x1) / (y2 - y1)
	x := x1

	for y := y1; y <= y2; y++ {
		(*intersections)[y] = append((*intersections)[y], int(x))
		x += m
	}
}

// DrawKochSnowflake draws a Koch snowflake.
func (ppm *PPM) DrawKochSnowflake(n int, start Point, width int, color Pixel) {
	ppm.drawKochSnowflakeSegment(n, start, Point{start.X + width, start.Y}, color)
	ppm.drawKochSnowflakeSegment(n, Point{start.X + width, start.Y}, Point{start.X + width/2, start.Y + int(math.Sqrt(3)*float64(width)/2)}, color)
	ppm.drawKochSnowflakeSegment(n, Point{start.X + width/2, start.Y + int(math.Sqrt(3)*float64(width)/2)}, start, color)
}

func (ppm *PPM) drawKochSnowflakeSegment(n int, p1, p2 Point, color Pixel) {
	if n == 0 {
		ppm.DrawLine(p1, p2, color)
		return
	}

	// Calculate one-third and two-thirds points of the segment
	oneThird := Point{
		X: (2*p1.X + p2.X) / 3,
		Y: (2*p1.Y + p2.Y) / 3,
	}
	twoThirds := Point{
		X: (p1.X + 2*p2.X) / 3,
		Y: (p1.Y + 2*p2.Y) / 3,
	}

	// Calculate equidistant point forming an equilateral triangle
	deltaX := twoThirds.X - oneThird.X
	deltaY := twoThirds.Y - oneThird.Y
	rotated := Point{
		X: oneThird.X + int(math.Cos(math.Pi/3)*float64(deltaX)-math.Sin(math.Pi/3)*float64(deltaY)),
		Y: oneThird.Y + int(math.Sin(math.Pi/3)*float64(deltaX)+math.Cos(math.Pi/3)*float64(deltaY)),
	}

	// Recursively draw the four segments of the Koch snowflake
	ppm.drawKochSnowflakeSegment(n-1, p1, oneThird, color)
	ppm.drawKochSnowflakeSegment(n-1, oneThird, rotated, color)
	ppm.drawKochSnowflakeSegment(n-1, rotated, twoThirds, color)
	ppm.drawKochSnowflakeSegment(n-1, twoThirds, p2, color)
}

// DrawSierpinskiTriangle draws a Sierpinski triangle.
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	ppm.drawSierpinskiTriangle(n, start, width, color)
}

func (ppm *PPM) drawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	if n == 0 {
		ppm.DrawFilledTriangle(
			start,
			Point{start.X + width, start.Y},
			Point{start.X + width/2, start.Y + int(math.Sqrt(3)*float64(width)/2)},
			color,
		)
		return
	}

	// Calculate midpoints of the three sides of the triangle
	mid1 := Point{(2*start.X + start.X + width) / 3, (2*start.Y + start.Y) / 3}
	mid2 := Point{(start.X + 2*start.X + width) / 3, (2*start.Y + start.Y) / 3}
	mid3 := Point{(start.X + start.X + width/2) / 2, (start.Y + start.Y + int(math.Sqrt(3)*float64(width)/2)) / 2}

	// Recursively draw the three sub-triangles
	ppm.drawSierpinskiTriangle(n-1, start, width/3, color)
	ppm.drawSierpinskiTriangle(n-1, mid1, width/3, color)
	ppm.drawSierpinskiTriangle(n-1, mid2, width/3, color)
	ppm.drawSierpinskiTriangle(n-1, mid3, width/3, color)
}

// NewPPM creates a new PPM image with the specified width and height.
func NewPPM(width, height int) *PPM {
	data := make([][]Pixel, height)
	for i := range data {
		data[i] = make([]Pixel, width)
	}
	return &PPM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: "P3",
		max:         255,
	}
}

// SavePNG saves the PPM image as a PNG file.
func (ppm *PPM) SavePNG(filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, ppm.width, ppm.height))

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			img.Set(x, y, color.RGBA{ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B, 255})
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
