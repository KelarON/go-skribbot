package ui

//go:generate fyne bundle --package resources --name ResourceImagePlaceholderPng -o resources/image_placeholder.go resources/image_placeholder.png

import (
	"go-skribbot/config"
	"go-skribbot/resources"
	"go-skribbot/ui/window"
	"go-skribbot/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	hook "github.com/robotn/gohook"
)

type UI struct {
	app        fyne.App
	mainWindow *window.MainWindow
	cfg        *config.Config
	drawer     *utils.Drawer
	logger     *utils.Logger
	searcher   *utils.Searcher
}

const APP_ID = "com.kelaron.goskribbot"

func ShowUI(cfg *config.Config, logger *utils.Logger, drawer *utils.Drawer, searcher *utils.Searcher) {

	ui := &UI{
		cfg:      cfg,
		drawer:   drawer,
		logger:   logger,
		searcher: searcher,
	}

	app := app.NewWithID(APP_ID)
	app.SetIcon(resources.ResourceIconPng)
	ui.app = app

	ui.ShowMainWindow()

	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()

	ui.app.Run()
}

func (ui *UI) ShowMainWindow() {
	ui.mainWindow = window.NewMainWindow(ui.app, ui.cfg, ui.searcher)
	ui.mainWindow.SetMaster()
	ui.mainWindow.SetOnClosed(ui.closeApp)
	ui.mainWindow.OnTableCellClick(ui.tableCellClick)
	ui.mainWindow.Show()
}

func (ui *UI) tableCellClick(id widget.TableCellID) {

	ui.mainWindow.ImageTable.Unselect(id)

	if id.Row < 0 || id.Col < 0 {
		return
	}

	matrix, err := ui.mainWindow.GetMatrixFromTable(id.Row, id.Col)
	if err != nil {
		ui.logger.Errorf("error finding image: %v", err)
		return
	}

	err = ui.drawer.DrawImage(&matrix)
	if err != nil {
		ui.logger.Errorf("error drawing image: %v", err)
		return
	}
}

func (ui *UI) closeApp() {
	ui.cfg.Save()
}
