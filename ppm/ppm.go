package ppm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

// PPM représente une image au format PPM.
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// ReadPPM lit une image PPM à partir d'un fichier et renvoie une structure PPM.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	header, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Vérifier le numéro magique pour déterminer le format P3 ou P6.
	magicNumber := strings.TrimSpace(header)
	if magicNumber != "P3" && magicNumber != "P6" {
		return nil, fmt.Errorf("format de fichier non supporté: %s", magicNumber)
	}

	// Lire les dimensions et la valeur maximale de couleur.
	var width, height, max int
	_, err = fmt.Fscanf(reader, "%d %d\n%d\n", &width, &height, &max)
	if err != nil {
		return nil, err
	}

	// Lire les données de l'image.
	data := make([][]Pixel, height)
	for i := range data {
		data[i] = make([]Pixel, width)
	}

	if magicNumber == "P3" {
		// Format P3 - ASCII.
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				var r, g, b int
				_, err = fmt.Fscanf(reader, "%d %d %d", &r, &g, &b)
				if err != nil {
					return nil, err
				}
				data[y][x] = Pixel{uint8(r), uint8(g), uint8(b)}
			}
		}
	} else {
		// Format P6 - Binaire.
		reader.ReadByte() // Consommer le caractère de nouvelle ligne après la valeur max.
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, _ := reader.ReadByte()
				g, _ := reader.ReadByte()
				b, _ := reader.ReadByte()
				data[y][x] = Pixel{r, g, b}
			}
		}
	}

	return &PPM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         uint8(max),
	}, nil
}

func main() {
	ppm, err := ReadPPM("PPM.ppm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier PPM :", err)
		return
	}
	fmt.Println("Image PPM lue avec succès :", ppm)
}
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}
func (ppm *PPM) At(x, y int) Pixel {
	if x < 0 || x >= ppm.width || y < 0 || y >= ppm.height {
		return Pixel{}
	}
	return ppm.data[y][x]
}
func (ppm *PPM) Set(x, y int, value Pixel) {
	if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
		ppm.data[y][x] = value
	}
}
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écrire l'en-tête du fichier PPM
	_, err = file.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", ppm.width, ppm.height))
	if err != nil {
		return err
	}

	// Écrire les données des pixels
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.At(x, y)
			_, err = file.WriteString(fmt.Sprintf("%d %d %d ", pixel.Red, pixel.Green, pixel.Blue))
			if err != nil {
				return err
			}
		}
		_, err = file.WriteString("\n") // Nouvelle ligne après chaque ligne de pixels
		if err != nil {
			return err
		}
	}

	return nil
}
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.At(x, y)
			// Inverser les couleurs en soustrayant la valeur de couleur actuelle de la valeur maximale
			pixel.Red = uint8(ppm.max) - pixel.Red
			pixel.Green = uint8(ppm.max) - pixel.Green
			pixel.Blue = uint8(ppm.max) - pixel.Blue
			ppm.Set(x, y, pixel)
		}
	}
}
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			// Échanger les pixels de gauche à droite
			oppositeX := ppm.width - 1 - x
			leftPixel := ppm.At(x, y)
			rightPixel := ppm.At(oppositeX, y)
			ppm.Set(x, y, rightPixel)
			ppm.Set(oppositeX, y, leftPixel)
		}
	}
}
func (ppm *PPM) Flop() {
	for x := 0; x < ppm.width; x++ {
		for y := 0; y < ppm.height/2; y++ {
			// Échanger les pixels du haut et du bas
			oppositeY := ppm.height - 1 - y
			topPixel := ppm.At(x, y)
			bottomPixel := ppm.At(x, oppositeY)
			ppm.Set(x, y, bottomPixel)
			ppm.Set(x, oppositeY, topPixel)
		}
	}
}
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	// Vérifier si le numéro magique est valide (P3 ou P6)
	if magicNumber != "P3" && magicNumber != "P6" {
		fmt.Println("Numéro magique non valide :", magicNumber)
		return
	}
	ppm.magicNumber = magicNumber
}
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	if maxValue > 255 {
		// Gérer l'erreur si la valeur maximale est supérieure à 255
		fmt.Println("La valeur maximale doit être inférieure ou égale à 255")
		return
	}
	ppm.max = uint8(maxValue)
}
func (ppm *PPM) ToPGM() {
	// Créer une nouvelle matrice pour les pixels PGM
	pgmData := make([][]uint8, ppm.height)
	for i := range pgmData {
		pgmData[i] = make([]uint8, ppm.width)
	}

	// Convertir chaque pixel PPM en une valeur de gris pour PGM
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			r := ppm.data[y][x].Red
			g := ppm.data[y][x].Green
			b := ppm.data[y][x].Blue
			// Calculer la luminance en utilisant une formule standard
			luminance := uint8(0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b))
			pgmData[y][x] = luminance
		}
	}

	// PGM représente une image au format PGM (Portable Graymap).
	type PGM struct {
		Width    int
		Height   int
		MaxValue uint8
		Pixels   [][]uint8
	}
}
func (ppm *PPM) ToPBM() *PBM {
	// Créer une nouvelle matrice pour les pixels PBM
	pbmData := make([][]bool, ppm.height)
	for i := range pbmData {
		pbmData[i] = make([]bool, ppm.width)
	}

	// Seuil pour déterminer si un pixel est blanc (true) ou noir (false)
	threshold := uint8(ppm.max / 2)

	// Convertir chaque pixel PPM en noir ou blanc pour PBM
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			r := ppm.data[y][x].Red
			g := ppm.data[y][x].Green
			b := ppm.data[y][x].Blue
			// Calculer la luminance
			luminance := uint8(0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b))
			// Appliquer le seuil pour déterminer la couleur
			pbmData[y][x] = luminance > threshold
		}
	}

	// Créer et retourner la nouvelle image PBM
	return &PBM{
		Width:  ppm.width,
		Height: ppm.height,
		Pixels: pbmData,
	}
}

// PBM représente une image au format PBM.
type PBM struct {
	Width, Height int
	Pixels        [][]bool
}
type Point struct {
	X, Y int
}

// NewPoint crée et retourne un nouveau point avec les coordonnées spécifiées.
func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y}
}

// Move déplace le point par un certain décalage en X et en Y.
func (p *Point) Move(dx, dy int) {
	p.X += dx
	p.Y += dy
}

// Distance calcule la distance entre deux points.
func (p *Point) Distance(other *Point) float64 {
	return math.Sqrt(float64((p.X-other.X)*(p.X-other.X) + (p.Y-other.Y)*(p.Y-other.Y)))
}
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	x1, y1 := p1.X, p1.Y
	x2, y2 := p2.X, p2.Y
	dx := abs(x2 - x1)
	dy := -abs(y2 - y1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx + dy

	for {
		ppm.Set(x1, y1, color)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}
		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
}

// abs renvoie la valeur absolue d'un entier.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	// Dessiner les côtés horizontaux du rectangle
	for x := p1.X; x < p1.X+width; x++ {
		ppm.Set(x, p1.Y, color)          // Côté supérieur
		ppm.Set(x, p1.Y+height-1, color) // Côté inférieur
	}
	// Dessiner les côtés verticaux du rectangle
	for y := p1.Y; y < p1.Y+height; y++ {
		ppm.Set(p1.X, y, color)         // Côté gauche
		ppm.Set(p1.X+width-1, y, color) // Côté droit
	}
}
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	// Remplir le rectangle avec la couleur spécifiée
	for y := p1.Y; y < p1.Y+height; y++ {
		for x := p1.X; x < p1.X+width; x++ {
			ppm.Set(x, y, color)
		}
	}
}
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	x, y, dx, dy := radius-1, 0, 1, 1
	err := dx - (radius * 2)

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
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (radius * 2)
		}
	}
}
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	// Tracer les trois côtés du triangle
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Fonction d'aide pour trouver le minimum et le maximum de deux valeurs
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// Fonction d'aide pour trier les points par ordre croissant de Y
	sortPointsByY := func(p1, p2, p3 Point) (Point, Point, Point) {
		if p1.Y > p2.Y {
			p1, p2 = p2, p1
		}
		if p1.Y > p3.Y {
			p1, p3 = p3, p1
		}
		if p2.Y > p3.Y {
			p2, p3 = p3, p2
		}
		return p1, p2, p3
	}

	// Trier les points pour simplifier le remplissage
	p1, p2, p3 = sortPointsByY(p1, p2, p3)

	// Fonction d'aide pour interpoler les valeurs de X le long des côtés du triangle
	interpolate := func(y, y1, y2, x1, x2 int) int {
		if y1 == y2 {
			return x1
		}
		return x1 + (x2-x1)*(y-y1)/(y2-y1)
	}

	// Remplir le triangle ligne par ligne
	for y := p1.Y; y <= p3.Y; y++ {
		var xStart, xEnd int
		if y < p2.Y {
			xStart = interpolate(y, p1.Y, p2.Y, p1.X, p2.X)
			xEnd = interpolate(y, p1.Y, p3.Y, p1.X, p3.X)
		} else {
			xStart = interpolate(y, p2.Y, p3.Y, p2.X, p3.X)
			xEnd = interpolate(y, p1.Y, p3.Y, p1.X, p3.X)
		}
		for x := min(xStart, xEnd); x <= max(xStart, xEnd); x++ {
			ppm.Set(x, y, color)
		}
	}
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	n := len(points)
	if n < 3 {
		// Pas assez de points pour former un polygone
		return
	}
	// Tracer des lignes entre les points consécutifs
	for i := 0; i < n-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	// Fermer le polygone en traçant une ligne entre le dernier point et le premier
	ppm.DrawLine(points[n-1], points[0], color)
}
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Fonction d'aide pour trouver le minimum et le maximum en Y
	minY, maxY := points[0].Y, points[0].Y
	for _, p := range points {
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Fonction d'aide pour trouver les intersections de l'arête avec une ligne horizontale
	getIntersections := func(y int) []int {
		var intersections []int
		n := len(points)
		for i := 0; i < n; i++ {
			p1 := points[i]
			p2 := points[(i+1)%n]

			if p1.Y == p2.Y { // Ignorer les arêtes horizontales
				continue
			}

			if (p1.Y <= y && p2.Y > y) || (p1.Y > y && p2.Y <= y) {
				// Trouver le point d'intersection avec la ligne horizontale
				x := p1.X + (y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y)
				intersections = append(intersections, x)
			}
		}
		sort.Ints(intersections) // Trier les intersections par ordre croissant en X
		return intersections
	}

	// Remplir le polygone ligne par ligne
	for y := minY; y <= maxY; y++ {
		intersections := getIntersections(y)
		for i := 0; i < len(intersections); i += 2 {
			if i+1 < len(intersections) {
				for x := intersections[i]; x <= intersections[i+1]; x++ {
					ppm.Set(x, y, color)
				}
			}
		}
	}
}

// Fonction d'aide pour calculer un point à une certaine distance et angle par rapport à un autre point
func pointAtDistanceAndAngle(from Point, distance float64, angle float64) Point {
	return Point{
		X: int(float64(from.X) + distance*math.Cos(angle)),
		Y: int(float64(from.Y) + distance*math.Sin(angle)),
	}
}

// Fonction récursive pour dessiner une courbe de Koch
func (ppm *PPM) drawKochCurve(p1, p2 Point, depth int, color Pixel) {
	if depth == 0 {
		ppm.DrawLine(p1, p2, color)
		return
	}

	// Diviser le segment en trois parties égales
	dx := float64(p2.X-p1.X) / 3.0
	dy := float64(p2.Y-p1.Y) / 3.0

	// Calculer les trois points de division
	pa := Point{p1.X + int(dx), p1.Y + int(dy)}
	pb := Point{p1.X + int(2*dx), p1.Y + int(2*dy)}

	// Calculer le point du pic
	angle := math.Atan2(dy, dx) - math.Pi/3
	pc := pointAtDistanceAndAngle(pa, math.Sqrt(dx*dx+dy*dy), angle)

	// Dessiner les quatre segments de la courbe de Koch
	ppm.drawKochCurve(p1, pa, depth-1, color)
	ppm.drawKochCurve(pa, pc, depth-1, color)
	ppm.drawKochCurve(pc, pb, depth-1, color)
	ppm.drawKochCurve(pb, p2, depth-1, color)
}

// DrawKochSnowflake dessine un flocon de neige de Koch.
func (ppm *PPM) DrawKochSnowflake(center Point, radius int, color Pixel) {
	// Calculer les trois sommets du triangle équilatéral initial
	p1 := pointAtDistanceAndAngle(center, float64(radius), -math.Pi/2)
	p2 := pointAtDistanceAndAngle(center, float64(radius), -math.Pi/2+2*math.Pi/3)
	p3 := pointAtDistanceAndAngle(center, float64(radius), -math.Pi/2+4*math.Pi/3)

	// Dessiner les trois côtés du flocon de neige de Koch
	depth := 4 // Profondeur de récursion, ajustez selon la complexité souhaitée
	ppm.drawKochCurve(p1, p2, depth, color)
	ppm.drawKochCurve(p2, p3, depth, color)
	ppm.drawKochCurve(p3, p1, depth, color)
}
