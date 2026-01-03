package main

import (
	"bufio"
	"fmt"
	"go-skribbot/config"
	"go-skribbot/model"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/nfnt/resize"
)

const (
	PICTURE_SIZE = 80
	PIXEL_SIZE   = 7
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		if err != config.ErrorNoConfigFile {
			log.Printf("error loading config: %v", err)
			exitMessage()
			return
		} else {
			cfg.Save()
		}
	}
	if cfg.PrintCoordsMode {
		time.Sleep(5 * time.Second)
		x, y := robotgo.Location()
		log.Printf("Your coordinates X: %v, Y: %v", x, y)
		exitMessage()
		return
	}

	startPositionX := cfg.PositionX
	startPositionY := cfg.PositionY - 600
	drawingType := cfg.DrawingType

	imageName := findImage()
	if imageName == "" {
		log.Printf("No image found in the current directory")
		exitMessage()
		return
	}

	file, err := os.Open(imageName)
	if err != nil {
		log.Printf("error opening file: %v", err)
		exitMessage()
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("error decoding image: %v", err)
		exitMessage()
		return
	}
	img = resize.Resize(PICTURE_SIZE, PICTURE_SIZE, img, resize.Bicubic)
	var matrix [PICTURE_SIZE][PICTURE_SIZE]*model.Color
	allowedColors := model.GetAllowedColors(cfg.PositionX, cfg.PositionY)

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

	switch drawingType {
	case model.DRAWING_TYPE_LINE:
		prepareBrush(startPositionX, startPositionY)
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
						robotgo.Move(startPositionX+x*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
						pX, pY := robotgo.Location()
						if pX != startPositionX+x*PIXEL_SIZE || pY != startPositionY+y*PIXEL_SIZE {
							robotgo.Move(startPositionX+x*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
						}
						robotgo.MouseDown(robotgo.Mleft)
						continue
					}
					if isLine {
						if isPont {
							robotgo.MouseUp(robotgo.Mleft)
						} else {
							robotgo.Move(startPositionX+(x-1)*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
							pX, pY := robotgo.Location()
							if pX < startPositionX-10 {
								log.Println("Stopped")
								exitMessage()
								return
							}
							if pX != startPositionX+(x-1)*PIXEL_SIZE || pY != startPositionY+y*PIXEL_SIZE {
								robotgo.Move(startPositionX+(x-1)*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
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
						robotgo.Move(startPositionX+PICTURE_SIZE*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
						pX, pY := robotgo.Location()
						if pX != startPositionX+PICTURE_SIZE*PIXEL_SIZE || pY != startPositionY+y*PIXEL_SIZE {
							robotgo.Move(startPositionX+PICTURE_SIZE*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
						}
						robotgo.MouseUp(robotgo.Mleft)
					}
				}
			}
		}
	case model.DRAWING_TYPE_POINT:
		prepareBrush(startPositionX, startPositionY)
		robotgo.MouseSleep = 3
		for _, color := range allowedColors[1:] {
			robotgo.Move(color.X, color.Y)
			robotgo.MilliSleep(100)
			robotgo.Click()
			for y, row := range matrix {
				for x, pixel := range row {
					if pixel.Id == color.Id {
						robotgo.Move(startPositionX+x*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
						pX, pY := robotgo.Location()
						if pX < startPositionX-10 {
							log.Println("Stopped")
							exitMessage()
							return
						}
						if pX != startPositionX+x*PIXEL_SIZE || pY != startPositionY+y*PIXEL_SIZE {
							robotgo.Move(startPositionX+x*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
						}
						robotgo.Click()
					}
				}
			}
		}
	default:
		log.Println("unknown drawing type, check drawing_type in config.yaml")
		exitMessage()
		return
	}
	log.Println("Success")
	exitMessage()
}

func findImage() string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if !entry.IsDir() && isImageFile(entry.Name()) {
			return entry.Name()
		}
	}
	return ""
}

// isImageFile checks if the file has a common image extension
func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".bmp", ".webp", ".svg":
		return true
	default:
		return false
	}
}

// set up starting brush and canvas state
func prepareBrush(startPositionX, startPositionY int) {
	time.Sleep(3 * time.Second)
	robotgo.Move(startPositionX, startPositionY)
	robotgo.KeyPress(robotgo.KeyC)
	robotgo.KeyPress(robotgo.KeyB)
	robotgo.ScrollSmooth(-1, 30, 20, 0)
	robotgo.ScrollSmooth(1, 3, 20, 0)
}

// exitMessage prompts the user to press Enter to exit the program
func exitMessage() {
	fmt.Println("Press Enter to exit")
	consoleReader := bufio.NewReaderSize(os.Stdin, 1)
	consoleReader.ReadByte()
}
