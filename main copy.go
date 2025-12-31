package main

// import (
// 	"drawbot/model"
// 	"fmt"
// 	"image"
// 	"os"
// 	"time"

// 	"github.com/go-vgo/robotgo"
// 	"github.com/nfnt/resize"
// )

// const (
// 	START_POSITION_X = model.CORE_POSITION_X       //818
// 	START_POSITION_Y = model.CORE_POSITION_Y - 600 //275
// 	PICTURE_SIZE     = 80
// 	PIXEL_SIZE       = 7
// )

// func main() {

// 	file, _ := os.Open("image.png")
// 	img, _, _ := image.Decode(file)
// 	img = resize.Resize(PICTURE_SIZE, PICTURE_SIZE, img, resize.Bicubic)
// 	var matrix [PICTURE_SIZE][PICTURE_SIZE]*model.Color
// 	allowedColors := model.GetAllowedColors()

// 	time.Sleep(3 * time.Second)

// 	fmt.Println(robotgo.Location())

// 	robotgo.MouseSleep = 3

// 	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
// 		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
// 			r, g, b, alf := img.At(x, y).RGBA()

// 			if alf < 5000 {
// 				r, g, b, alf = 65535, 65535, 65535, 65535
// 			}

// 			curCol := model.Color{
// 				R: r,
// 				G: g,
// 				B: b,
// 				X: x,
// 				Y: y,
// 			}
// 			matrix[y][x] = curCol.FindClosest(allowedColors)
// 		}
// 	}

// 	left, right := 0, 0

// 	for _, color := range allowedColors[1:] {
// 		if left == 0 {
// 			left = color.Id
// 			robotgo.Move(color.X, color.Y)
// 			robotgo.MilliSleep(100)
// 			robotgo.Click("left")
// 			continue
// 		} else {
// 			right = color.Id
// 			robotgo.Move(color.X, color.Y)
// 			robotgo.MilliSleep(100)
// 			robotgo.Click("right")
// 		}
// 		for y, row := range matrix {
// 			for x, pixel := range row {
// 				if pixel.Id == left {
// 					robotgo.Move(START_POSITION_X+x*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
// 					pX, pY := robotgo.Location()
// 					if pX != START_POSITION_X+x*PIXEL_SIZE || pY != START_POSITION_Y+y*PIXEL_SIZE {
// 						robotgo.Move(START_POSITION_X+x*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
// 					}
// 					robotgo.Click()
// 				} else {
// 					if pixel.Id == right {
// 						robotgo.Move(START_POSITION_X+x*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
// 						pX, pY := robotgo.Location()
// 						if pX != START_POSITION_X+x*PIXEL_SIZE || pY != START_POSITION_Y+y*PIXEL_SIZE {
// 							robotgo.Move(START_POSITION_X+x*PIXEL_SIZE, START_POSITION_Y+y*PIXEL_SIZE)
// 						}
// 						robotgo.Click("right")
// 					}
// 				}
// 			}
// 		}
// 		left, right = 0, 0
// 	}
// }
