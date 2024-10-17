package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
)

func main() {
	if err := resizeThumbnails("assets"); err != nil {
		log.Print(err)
	}
}

func resizeThumbnails(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	var errorsSlice []error
	for _, file := range files {
		if isThumbnail(file.Name()) {
			inputPath := filepath.Join(dir, file.Name())
			outputPath := filepath.Join(dir, "resized_"+file.Name())

			if err = processImage(inputPath, outputPath, file.Name()); err != nil {
				errorsSlice = append(errorsSlice, fmt.Errorf("failed to resize %s: %v", file.Name(), err))
			}
		}
	}
	return errors.Join(errorsSlice...)
}

func isThumbnail(fileName string) bool {
	return strings.HasPrefix(fileName, "thumbnail_")
}

func processImage(inputPath, outputPath, fileName string) error {
	switch {
	case strings.HasSuffix(fileName, ".jpeg") || strings.HasSuffix(fileName, ".jpg"):
		return resizeJpgToSquare(inputPath, outputPath)
	case strings.HasSuffix(fileName, ".webp"):
		return resizeWebpToSquare(inputPath, outputPath)
	default:
		return nil
	}
}

func resizeJpgToSquare(inputPath, outputPath string) error {
	return resizeImage(inputPath, outputPath, jpegEncodeWrapper)
}

func resizeWebpToSquare(inputPath, outputPath string) error {
	return resizeImage(inputPath, outputPath, func(w io.Writer, img image.Image) error {
		return webp.Encode(w, img, &webp.Options{Lossless: true})
	})
}

func jpegEncodeWrapper(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, nil)
}

func resizeImage(inputPath, outputPath string, encodeFunc func(io.Writer, image.Image) error) error {
	imgFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	squareImg := createSquareImage(img)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return encodeFunc(outFile, squareImg)
}

func createSquareImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newSize := max(width, height)
	squareImg := image.NewRGBA(image.Rect(0, 0, newSize, newSize))
	draw.Draw(squareImg, squareImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	offsetX := (newSize - width) / 2
	offsetY := (newSize - height) / 2

	draw.Draw(squareImg,
		image.Rect(offsetX, offsetY, offsetX+width, offsetY+height),
		img,
		bounds.Min,
		draw.Over)

	return squareImg
}
