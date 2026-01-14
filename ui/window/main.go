package window

import (
	"errors"
	"go-skribbot/config"
	"go-skribbot/model"
	"go-skribbot/resources"
	"go-skribbot/utils"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/sqweek/dialog"
)

type MainWindow struct {
	fyne.Window

	searcher      *utils.Searcher
	searchResults []SearchResult
	searchButton  *widget.Button
	ImageTable    *widget.Table

	waitingCoords bool
	coordsChan    chan Coords
}

type SearchResult struct {
	image  *image.RGBA
	matrix [utils.PICTURE_SIZE][utils.PICTURE_SIZE]*model.Color
}

type Coords struct {
	X int
	Y int
}

const (
	DEFAULT_WINDOW_WIDTH  = 1000
	DEFAULT_WINDOW_HEIGHT = 720
)

func NewMainWindow(app fyne.App, cfg *config.Config, searcher *utils.Searcher) *MainWindow {

	mw := &MainWindow{
		searcher:   searcher,
		coordsChan: make(chan Coords),
	}

	mw.startCoordsFinder()

	mw.Window = app.NewWindow("Go-Skribbot")
	mw.SetFixedSize(true)
	mw.Resize(fyne.NewSize(DEFAULT_WINDOW_WIDTH, DEFAULT_WINDOW_HEIGHT))

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		NewSettingsWindow(app, cfg, &mw.waitingCoords, mw.coordsChan).Show()
	})
	settingsButton.Resize(fyne.NewSize(40, 40))
	settingsButton.Move(fyne.NewPos(940, 10))

	imageTable := widget.NewTable(
		func() (rows int, cols int) { return 2, 3 },
		func() fyne.CanvasObject {
			placeholder := canvas.NewImageFromResource(resources.ResourceImagePlaceholderPng)
			placeholder.ScaleMode = canvas.ImageScalePixels
			placeholder.SetMinSize(fyne.NewSize(320, 320))
			image := canvas.NewImageFromImage(nil)
			image.ScaleMode = canvas.ImageScalePixels
			image.SetMinSize(fyne.NewSize(320, 320))
			image.Hide()
			startIcon := widget.NewIcon(theme.ConfirmIcon())
			startIcon.Resize(fyne.NewSize(320, 320))
			startIcon.Theme().Color(theme.ColorGreen, theme.VariantDark)
			startIcon.Hide()
			return container.NewStack(placeholder, image, startIcon)
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			container := co.(*fyne.Container)
			placeholder := container.Objects[0].(*canvas.Image)
			image := container.Objects[1].(*canvas.Image)
			startIcon := container.Objects[2].(*widget.Icon)
			if len(mw.searchResults) > 0 {
				image.Image = mw.searchResults[tci.Row*3+tci.Col].image
				image.Refresh()
				placeholder.Hide()
				image.Show()
			} else {
				image.Hide()
				startIcon.Hide()
				placeholder.Show()
			}
		},
	)
	mw.ImageTable = imageTable
	imageTable.Resize(fyne.NewSize(970, 650))
	imageTable.Move(fyne.NewPos(10, 60))

	searchEntry := widget.NewEntry()
	searchEntry.Move(fyne.NewPos(10, 10))
	searchEntry.Resize(fyne.NewSize(320, 40))
	searchEntry.PlaceHolder = "Enter image search query"
	searchEntry.OnSubmitted = func(s string) {
		if cfg.PositionX == 0 && cfg.PositionY == 0 {
			go dialog.Message("Set the coordinates in the settings to draw").Title("Warning").Info()
		}
		fyne.Do(mw.searchButton.Disable)
		mw.searchResults = make([]SearchResult, 0, utils.IMAGE_COUNT)
		images, err := mw.searcher.SearchImages(searchEntry.Text)
		if err != nil {
			mw.searchResults = nil
		}
		for _, img := range images {
			matrix, image, err := utils.PrepareImage(img)
			if err != nil {
				mw.searchResults = nil
			}
			mw.searchResults = append(mw.searchResults, SearchResult{image: image, matrix: matrix})
		}
		fyne.Do(imageTable.Refresh)
		fyne.Do(mw.searchButton.Enable)
	}
	searchEntry.Refresh()

	searchButton := widget.NewButton("Search", func() {
		if cfg.PositionX == 0 && cfg.PositionY == 0 {
			go dialog.Message("Set the coordinates in the settings to draw").Title("Warning").Info()
		}
		fyne.Do(mw.searchButton.Disable)
		mw.searchResults = make([]SearchResult, 0, utils.IMAGE_COUNT)
		images, err := mw.searcher.SearchImages(searchEntry.Text)
		if err != nil {
			mw.searchResults = nil
		}
		for _, img := range images {
			matrix, image, err := utils.PrepareImage(img)
			if err != nil {
				mw.searchResults = nil
			}
			mw.searchResults = append(mw.searchResults, SearchResult{image: image, matrix: matrix})
		}
		fyne.Do(imageTable.Refresh)
		fyne.Do(mw.searchButton.Enable)
	})
	mw.searchButton = searchButton
	searchButton.Resize(fyne.NewSize(100, 40))
	searchButton.Move(fyne.NewPos(340, 10))

	mw.SetContent(container.NewWithoutLayout(
		searchEntry,
		searchButton,
		settingsButton,
		imageTable,
	))

	return mw
}

func (mw *MainWindow) OnTableCellClick(f func(id widget.TableCellID)) {
	mw.ImageTable.OnSelected = f
}

func (mw *MainWindow) GetMatrixFromTable(row, col int) ([utils.PICTURE_SIZE][utils.PICTURE_SIZE]*model.Color, error) {

	if len(mw.searchResults) == 0 {
		return [utils.PICTURE_SIZE][utils.PICTURE_SIZE]*model.Color{}, errors.New("find images first")
	}

	if row*3+col >= utils.IMAGE_COUNT {
		return mw.searchResults[0].matrix, nil
	} else {
		return mw.searchResults[row*3+col].matrix, nil
	}
}

func (mw *MainWindow) startCoordsFinder() {
	hook.Register(hook.MouseDown, nil, func(e hook.Event) {
		if e.Button == hook.MouseMap["right"] {
			if mw.waitingCoords {
				x, y := robotgo.Location()
				mw.coordsChan <- Coords{X: x, Y: y}
			}
		}
	})
}
