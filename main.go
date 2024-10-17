package main

import (
	"log"

	"github.com/ValeryVerkhoturov/image-resize/image_resize"
)

func main() {
	if err := image_resize.ResizeThumbnails("assets"); err != nil {
		log.Print(err)
	}
}
