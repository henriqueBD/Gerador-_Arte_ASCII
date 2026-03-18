package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	const path = "."
	const format = "png"
	const nameTiles = "chars"
	const tileSide = 8

	var nameImg string
	fmt.Print("> ")
	fmt.Scanln(&nameImg)

	tiles, numTiles := openAsciiTile(nameTiles, path, "png", tileSide)
	if len(tiles) == 0 {
		fmt.Println("Error: No tiles loaded.")
		return
	}

	img := openImage(nameImg, path, format)
	if img == nil {
		return
	}

	grayImg := convertGrayscale(img)
	grayImg = downscaleImage(grayImg, tileSide)

	arr := quantizeImage(grayImg, numTiles)

	res := turnIntoAscii(arr, tiles)

	if res != nil {
		savePngImage(res, path, nameImg)
	}
}

func openImage(name, path, formatRequired string) image.Image {
	s := path + "/images/" + name + "." + formatRequired
	fl, err := os.Open(s)
	if err != nil {
		fmt.Println("Error in opening image 1")
		fmt.Println(err)
		return nil
	}
	defer fl.Close()

	img, format, err := image.Decode(fl)
	if err != nil {
		fmt.Println("Error in decoding image")
		fmt.Println(err)
		return nil
	}

	if format != formatRequired {
		fmt.Println("Format must be " + formatRequired)
		return nil
	}

	return img
}

func convertGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	newImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayPixel := color.GrayModel.Convert(img.At(x, y))
			newImg.Set(x, y, grayPixel)
		}
	}

	return newImg
}

func savePngImage(img image.Image, path, name string) {
	path += "/images/" + name + "_ASCII" + ".png"

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Println("Error encoding image")
		fmt.Println(err)
		return
	}

	fmt.Println("Image converted successfully!")
}

func downscaleImage(img *image.Gray, factor int) *image.Gray {
	if factor < 1 {
		return img
	}

	return nearestNeighbour(img, factor)
}

func openAsciiTile(name, path, formatRequired string, tileSide int) ([]image.Image, uint8) {
	s := path + "/" + name + "." + formatRequired
	fl, err := os.Open(s)
	if err != nil {
		fmt.Println("Error in opening tile image")
		fmt.Println(err)
		return nil, 0
	}
	defer fl.Close()

	img, format, err := image.Decode(fl)
	if err != nil {
		fmt.Println("Error in decoding tile image")
		fmt.Println(err)
		return nil, 0
	}
	if format != formatRequired {
		fmt.Println("Format must be " + formatRequired)
		return nil, 0
	}

	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	numTiles := w / tileSide

	res := make([]image.Image, 0, numTiles)

	for i := 0; i < numTiles; i++ {
		tileBounds := image.Rect(tileSide*i, 0, tileSide*(i+1), h)
		tile := image.NewGray(tileBounds)
		draw.Draw(tile, tileBounds, img, image.Point{X: i * tileSide, Y: 0}, draw.Src)
		res = append(res, tile)
	}

	return res, uint8(numTiles)
}

func quantizeImage(img *image.Gray, steps uint8) [][]uint8 {
	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	res := make([][]uint8, h)
	for i := range res {
		res[i] = make([]uint8, w)
	}

	if steps == 0 {
		return res
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.GrayAt(x, y)
			quantizedValue := uint8((uint16(originalPixel.Y) * uint16(steps-1)) / 255)
			res[y-bounds.Min.Y][x-bounds.Min.X] = quantizedValue
		}
	}

	return res
}

func turnIntoAsciiText(arr [][]uint8, path string) {
	chars := []uint8{' ', '.', ':', 'c', 'o', 'P', 'O', '?', 'Z', '@'}

	var line, res string

	for y := range arr {
		line = ""
		for x := range arr[y] {
			// Added safety check in case the quantized value exceeds char array bounds
			idx := arr[y][x]
			if int(idx) >= len(chars) {
				idx = uint8(len(chars) - 1)
			}
			line += string(chars[idx])
		}
		line += "\n"
		res += line
	}

	file, err := os.Create(path + "/res.txt")
	if err != nil {
		fmt.Println("Error in creating file:", err)
		return
	}
	defer file.Close()

	file.WriteString(res)
}

func turnIntoAscii(arr [][]uint8, tiles []image.Image) *image.Gray {
	if len(tiles) == 0 || len(arr) == 0 || len(arr[0]) == 0 {
		fmt.Println("Invalid array or tiles")
		return nil
	}

	bounds := tiles[0].Bounds()
	side := bounds.Max.X - bounds.Min.X

	if side == 0 {
		fmt.Println("Error in tile size")
		return nil
	}

	w := len(arr[0])
	h := len(arr)

	newImg := image.NewGray(image.Rect(0, 0, w*side, h*side))
	for y := range arr {
		for x := range arr[y] {
			rect := image.Rect(x*side, y*side, (x+1)*side, (y+1)*side)
			index := arr[y][x]

			// Safety check for out-of-bounds mapping
			if int(index) >= len(tiles) {
				index = uint8(len(tiles) - 1)
			}

			boundsTmp := tiles[index].Bounds()
			point := image.Point{X: boundsTmp.Min.X, Y: boundsTmp.Min.Y}
			draw.Draw(newImg, rect, tiles[index], point, draw.Src)
		}
	}

	return newImg
}

func nearestNeighbour(img *image.Gray, factor int) *image.Gray {
	bounds := img.Bounds()
	OldW := bounds.Max.X - bounds.Min.X
	OldH := bounds.Max.Y - bounds.Min.Y

	w := int(OldW / factor)
	h := int(OldH / factor)

	newBounds := image.Rect(0, 0, w, h)
	newImg := image.NewGray(newBounds)

	for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
		for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
			// Fixed coordinates: Needs to offset from the source image's bounds.Min
			srcX := bounds.Min.X + (x * factor)
			srcY := bounds.Min.Y + (y * factor)
			newImg.Set(x, y, img.At(srcX, srcY))
		}
	}

	return newImg
}
