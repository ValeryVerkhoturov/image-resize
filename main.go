package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

func main() {
	err := resizeToSquare(filepath.Join("assets", "input.jpg"), filepath.Join("assets", "output.jpg"))
	if err != nil {
		log.Print(err)
	}

	err = resizeToSquare(filepath.Join("assets", "input2.jpg"), filepath.Join("assets", "output2.jpg"))
	if err != nil {
		log.Print(err)
	}
}

func resizeToSquare(inputPath, outputPath string) error {
	imgFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var newSize int
	if width > height {
		newSize = width
	} else {
		newSize = height
	}

	squareImg := image.NewRGBA(image.Rect(0, 0, newSize, newSize))
	draw.Draw(squareImg, squareImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	offsetX := (newSize - width) / 2
	offsetY := (newSize - height) / 2

	draw.Draw(squareImg, image.Rect(offsetX, offsetY, offsetX+width, offsetY+height), img, bounds.Min, draw.Over)

	outFile, err := os.Create(outputPath)
	defer outFile.Close()
	if err != nil {
		return err
	}
	return jpeg.Encode(outFile, squareImg, nil)
}
