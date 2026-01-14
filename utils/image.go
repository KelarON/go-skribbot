package utils

import (
	"fmt"
	"go-skribbot/model"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

const PICTURE_SIZE = 80

func PrepareImage(img image.Image) ([PICTURE_SIZE][PICTURE_SIZE]*model.Color, *image.RGBA, error) {

	allowedColors := model.GetAllowedColors(0, 0)

	img, err := cropToSquareImage(img)
	if err != nil {
		return [PICTURE_SIZE][PICTURE_SIZE]*model.Color{}, nil, fmt.Errorf("error cropping image: %v", err)
	}

	img = resize.Resize(PICTURE_SIZE, PICTURE_SIZE, img, resize.Bicubic)
	var matrix [PICTURE_SIZE][PICTURE_SIZE]*model.Color

	allowedImg := image.NewRGBA(image.Rect(0, 0, 79, 79))

	offsetX := img.Bounds().Min.X
	offsetY := img.Bounds().Min.Y

	for y := 0; y < PICTURE_SIZE; y++ {
		for x := 0; x < PICTURE_SIZE; x++ {
			r, g, b, alf := img.At(x+offsetX, y+offsetY).RGBA()

			if alf < 5000 {
				r, g, b, alf = 65535, 65535, 65535, 65535
			}

			curCol := model.Color{
				R: r,
				G: g,
				B: b,
				X: x,
				Y: y,
			}
			matrix[y][x] = curCol.FindClosest(allowedColors)
			allowedImg.Set(x, y, color.RGBA64{
				R: uint16(matrix[y][x].R),
				G: uint16(matrix[y][x].G),
				B: uint16(matrix[y][x].B),
				A: 65535,
			})
		}
	}
	return matrix, allowedImg, nil
}

func cropToSquareImage(img image.Image) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	var x0, y0, x1, y1 int

	if img.Bounds().Size().X < img.Bounds().Size().Y {
		x0, y0 = 0, 0
		x1, y1 = img.Bounds().Size().X, img.Bounds().Size().X
	} else {
		if img.Bounds().Size().X > img.Bounds().Size().Y {
			y0 = 0
			offset := (img.Bounds().Size().X - img.Bounds().Size().Y) / 2
			x0 = offset
			y1 = img.Bounds().Size().Y
			x1 = img.Bounds().Size().Y + offset
		} else {
			return img, nil
		}
	}

	crop := image.Rect(x0, y0, x1, y1)

	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

func WriteImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
