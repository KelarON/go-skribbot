package utils

import (
	"errors"
	"go-skribbot/model"
	"sync/atomic"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

const PIXEL_SIZE = 7

type Drawer struct {
	drawingType *model.DrawingType
	positionX   *int
	positionY   *int
	busy        atomic.Bool
}

func NewDrawer(drawingType *model.DrawingType, positionX *int, positionY *int) *Drawer {
	drawer := &Drawer{
		drawingType: drawingType,
		positionX:   positionX,
		positionY:   positionY,
	}

	drawer.startStopper()

	return drawer
}

// prepareBrush sets correct brush size
func prepareBrush(startPositionX, startPositionY int) {
	robotgo.Move(startPositionX, startPositionY)
	robotgo.KeyPress(robotgo.KeyC)
	robotgo.KeyPress(robotgo.KeyB)
	robotgo.ScrollSmooth(-1, 30, 20, 0)
	robotgo.ScrollSmooth(1, 3, 20, 0)
}

// DrawImage draws the image on the screen based on the provided matrix.
func (d *Drawer) DrawImage(matrix *[PICTURE_SIZE][PICTURE_SIZE]*model.Color) error {

	if d.busy.Load() || (*d.positionX == 0 && *d.positionY == 0) {
		return nil
	}

	d.busy.Store(true)

	time.Sleep(3 * time.Second)

	// get required params for drawing
	startPositionX := *d.positionX
	startPositionY := *d.positionY - 600
	allowedColors := model.GetAllowedColors(*d.positionX, *d.positionY)
	usedColors := make([]*model.Color, 0)
	for _, color := range allowedColors[1:] {
		var found bool
		for _, row := range *matrix {
			for _, pixel := range row {
				if pixel.Id != 0 && pixel.Id == color.Id {
					usedColors = append(usedColors, color)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	prepareBrush(startPositionX, startPositionY)

	switch *d.drawingType {
	case model.DRAWING_TYPE_LINE:
		robotgo.MouseSleep = 20
		for _, color := range usedColors {
			robotgo.Move(color.X, color.Y)
			robotgo.MilliSleep(100)
			robotgo.Click()
			for y, row := range *matrix {
				isLine := false
				isPont := false
				for x, pixel := range row {
					if !d.busy.Load() {
						robotgo.MouseUp(robotgo.Mleft)
						return nil
					}
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
		robotgo.MouseSleep = 3
		for _, color := range usedColors {
			robotgo.Move(color.X, color.Y)
			robotgo.MilliSleep(100)
			robotgo.Click()
			for y, row := range *matrix {
				for x, pixel := range row {
					if !d.busy.Load() {
						robotgo.MouseUp(robotgo.Mleft)
						return nil
					}
					if pixel.Id == color.Id {
						robotgo.Move(startPositionX+x*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
						pX, pY := robotgo.Location()
						if pX != startPositionX+x*PIXEL_SIZE || pY != startPositionY+y*PIXEL_SIZE {
							robotgo.Move(startPositionX+x*PIXEL_SIZE, startPositionY+y*PIXEL_SIZE)
						}
						robotgo.Click()
					}
				}
			}
		}
	default:
		d.busy.Store(false)
		return errors.New("unknown drawing type, check drawing_type in config.yaml")
	}
	d.busy.Store(false)
	return nil
}

func (d *Drawer) startStopper() {
	hook.Register(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
		d.busy.Store(false)
	})
}
