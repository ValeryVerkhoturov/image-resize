package image_resize

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
)

func ResizeThumbnails(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	var errorsSlice []error
	for _, file := range files {
		if isThumbnail(file.Name()) {
			inputPath := filepath.Join(dir, file.Name())
			outputPath := filepath.Join(dir, "resized_"+file.Name())

			if err = ProcessImage(inputPath, outputPath, file.Name()); err != nil {
				errorsSlice = append(errorsSlice, fmt.Errorf("failed to resize %s: %v", file.Name(), err))
			}
		}
	}
	return errors.Join(errorsSlice...)
}

func isThumbnail(fileName string) bool {
	return strings.HasPrefix(fileName, "thumbnail_")
}

func ProcessImage(inputPath, outputPath, fileName string) error {
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
	return resizeImage(inputPath, outputPath, jpegEncodeWrapper, &color.White)
}

func jpegEncodeWrapper(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, nil)
}

func resizeWebpToSquare(inputPath, outputPath string) error {
	return resizeImage(inputPath, outputPath, webpEncodeWrapper, nil)
}

func webpEncodeWrapper(w io.Writer, img image.Image) error {
	return webp.Encode(w, img, &webp.Options{Lossless: true, Exact: true})
}

func resizeImage(inputPath,
	outputPath string,
	encodeFunc func(io.Writer, image.Image) error,
	bordersColor *color.Gray16) error {
	imgFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	squareImg := createSquareImage(img, bordersColor)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return encodeFunc(outFile, squareImg)
}

func createSquareImage(img image.Image, bordersColor *color.Gray16) *image.RGBA {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newSize := max(width, height)
	squareImg := image.NewRGBA(image.Rect(0, 0, newSize, newSize))

	if bordersColor != nil {
		draw.Draw(squareImg, squareImg.Bounds(), &image.Uniform{bordersColor}, image.Point{}, draw.Src)
	}

	offsetX := (newSize - width) / 2
	offsetY := (newSize - height) / 2

	draw.Draw(squareImg,
		image.Rect(offsetX, offsetY, offsetX+width, offsetY+height),
		img,
		bounds.Min,
		draw.Over)

	return squareImg
}
