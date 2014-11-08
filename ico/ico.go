package ico

import "image"

// ICO represents the possibly multiple images stored in a ICO file.
type ICO struct {
	Num   int           // Total number of images
	Image []image.Image // The images themselves
}
