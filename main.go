package main

import (
	"drawbot/model"
	"fmt"
	"image"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/nfnt/resize"
)

type drawingType int

const (
	DRAWING_TYPE_LINE = iota
	DRAWING_TYPE_POINT
)

const (
	START_POSITION_X = model.CORE_POSITION_X       //818
	START_POSITION_Y = model.CORE_POSITION_Y - 600 //275
	PICTURE_SIZE     = 80
	PIXEL_SIZE       = 7
	DRAWING_TYPE     = DRAWING_TYPE_LINE
)

func main() {

	file, _ := os.Open("image.png")
	img, _, _ := image.Decode(file)
	img = resize.Resize(PICTURE_SIZE, PICTURE_SIZE, img, resize.Bicubic)
	var matrix [PICTURE_SIZE][PICTURE_SIZE]*model.Color
	allowedColors := model.GetAllowedColors()

	time.Sleep(3 * time.Second)

	fmt.Println(robotgo.Location())

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, alf := img.At(x, y).RGBA()

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
		}
	}

	switch DRAWING_TYPE {
	case DRAWING_TYPE_LINE:
		robotgo.MouseSleep = 20
		for _, color := range allowedColors[1:] {
			robotgo.Move(color.X, color.Y)
			robotgo.MilliSleep(100)
			robotgo.Click()
			for y, row := range matrix {
				isLine := false
				isPont := false
				for x, pixel := range row {
					if pixel.Id == color.Id {
						if isLine {
							isPont = false
							continue
						}
						isLine = true
						isPont = true
						robotgo.Move(START_POSITION_X+x*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
						pX, pY := robotgo.Location()
						if pX != START_POSITION_X+x*PIXEL_SIZE || pY != START_POSITION_Y+y*PIXEL_SIZE {
							robotgo.Move(START_POSITION_X+x*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
						}
						robotgo.MouseDown(robotgo.Mleft)
						continue
					}
					if isLine {
						if isPont {
							robotgo.MouseUp(robotgo.Mleft)
						} else {
							robotgo.Move(START_POSITION_X+(x-1)*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
							pX, pY := robotgo.Location()
							if pX != START_POSITION_X+(x-1)*PIXEL_SIZE || pY != START_POSITION_Y+y*PIXEL_SIZE {
								robotgo.Move(START_POSITION_X+(x-1)*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
							}
							robotgo.MouseUp(robotgo.Mleft)
						}
					}
					isLine = false
				}
				if isLine {
					if isPont {
						robotgo.MouseUp(robotgo.Mleft)
					} else {
						robotgo.Move(START_POSITION_X+PICTURE_SIZE*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
						pX, pY := robotgo.Location()
						if pX != START_POSITION_X+PICTURE_SIZE*PIXEL_SIZE || pY != START_POSITION_Y+y*PIXEL_SIZE {
							robotgo.Move(START_POSITION_X+PICTURE_SIZE*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
						}
						robotgo.MouseUp(robotgo.Mleft)
					}
				}
			}
		}
	case DRAWING_TYPE_POINT:
		robotgo.MouseSleep = 3
		for _, color := range allowedColors[1:] {
			robotgo.Move(color.X, color.Y)
			robotgo.MilliSleep(100)
			robotgo.Click()
			for y, row := range matrix {
				for x, pixel := range row {
					if pixel.Id == color.Id {
						robotgo.Move(START_POSITION_X+x*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
						pX, pY := robotgo.Location()
						if pX != START_POSITION_X+x*PIXEL_SIZE || pY != START_POSITION_Y+y*PIXEL_SIZE {
							robotgo.Move(START_POSITION_X+x*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
						}
						robotgo.Click()
					}
				}
			}
		}
	}
}
