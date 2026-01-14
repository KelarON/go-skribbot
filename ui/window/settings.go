package window

import (
	"go-skribbot/config"
	"go-skribbot/model"
	"go-skribbot/resources"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

type SettingsWindow struct {
	fyne.Window

	cfg          *config.Config
	localCfg     config.Config
	coordsLabelX *widget.Label
	coordsLabelY *widget.Label

	waitingCoords *bool
	coordsChan    chan Coords
}

func NewSettingsWindow(app fyne.App, cfg *config.Config, wc *bool, coordsChan chan Coords) *SettingsWindow {

	sw := &SettingsWindow{
		cfg:           cfg,
		localCfg:      *cfg,
		waitingCoords: wc,
		coordsChan:    coordsChan,
	}

	settingsWindow := app.NewWindow("Settings")
	settingsWindow.SetFixedSize(true)
	settingsWindow.Resize(fyne.NewSize(300, 300))
	settingsWindow.CenterOnScreen()
	sw.Window = settingsWindow

	coordsLabel := widget.NewLabel("White color position       :")
	coordsLabel.Move(fyne.NewPos(10, 10))
	coordsLabelX := widget.NewLabel(strconv.Itoa(sw.localCfg.PositionX))
	sw.coordsLabelX = coordsLabelX
	coordsLabelX.Move(fyne.NewPos(10, 40))
	coordsLabelY := widget.NewLabel(strconv.Itoa(sw.localCfg.PositionY))
	sw.coordsLabelY = coordsLabelY
	coordsLabelY.Move(fyne.NewPos(60, 40))

	coordsHelpButton := widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
		hw := app.NewWindow("Help")
		hw.CenterOnScreen()
		img := canvas.NewImageFromResource(resources.ResourceCoordsHelpPng)
		img.FillMode = canvas.ImageFillOriginal
		hw.SetContent(img)
		hw.Show()
	})
	coordsHelpButton.Resize(fyne.NewSize(20, 20))
	coordsHelpButton.Move(fyne.NewPos(152, 18))

	coordsButton := widget.NewButton("Change", func() {
		go sw.setNewCoords(&sw.localCfg.PositionX, &sw.localCfg.PositionY)
	})
	coordsButton.Move(fyne.NewPos(150, 50))
	coordsButton.Resize(fyne.NewSize(70, 20))

	drawingTypeLabel := widget.NewLabel("Drawing type:")
	drawingTypeLabel.Move(fyne.NewPos(10, 80))

	drawingTypeRadioGroup := widget.NewRadioGroup(model.GetAllDrawingStatuses(), func(s string) {
		sw.localCfg.DrawingType = model.DrawingType(s)
	})
	drawingTypeRadioGroup.Move(fyne.NewPos(10, 110))
	drawingTypeRadioGroup.SetSelected(string(sw.localCfg.DrawingType))

	confirmButton := widget.NewButton("OK", func() {
		*sw.cfg = sw.localCfg
		sw.cfg.Save()
		sw.Close()
	})
	confirmButton.Move(fyne.NewPos(155, 250))
	confirmButton.Resize(fyne.NewSize(135, 40))

	cancelButton := widget.NewButton("Cancel", func() { sw.Close() })
	cancelButton.Move(fyne.NewPos(10, 250))
	cancelButton.Resize(fyne.NewSize(135, 40))

	sw.SetContent(container.NewWithoutLayout(
		coordsLabel,
		coordsHelpButton,
		coordsLabelX,
		coordsLabelY,
		coordsButton,
		drawingTypeLabel,
		drawingTypeRadioGroup,
		confirmButton,
		cancelButton,
	))

	sw.SetOnClosed(func() {
		*sw.waitingCoords = false
	})

	return sw
}

func (sw *SettingsWindow) setNewCoords(x *int, y *int) {
	if *sw.waitingCoords {
		return
	}
	*sw.waitingCoords = true
	fyne.Do(func() { sw.coordsLabelX.SetText("right") })
	fyne.Do(func() { sw.coordsLabelY.SetText("click") })
	coords := <-sw.coordsChan
	*sw.waitingCoords = false
	*x = coords.X
	*y = coords.Y
	fyne.Do(sw.RequestFocus)
	fyne.Do(func() { sw.coordsLabelX.SetText(strconv.Itoa(sw.localCfg.PositionX)) })
	fyne.Do(func() { sw.coordsLabelY.SetText(strconv.Itoa(sw.localCfg.PositionY)) })
	dialog.Message("Coordinates updated!").Title("Success").Info()
}
