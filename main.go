package main

import (
	"go-skribbot/config"
	"go-skribbot/ui"
	"go-skribbot/utils"

	hook "github.com/robotn/gohook"
)

func main() {

	logger := utils.NewLogger()
	defer logger.Close()

	cfg, err := config.LoadConfig()
	if err != nil {
		if err != config.ErrorNoConfigFile {
			logger.Errorf("error loading config: %v", err)
			return
		} else {
			cfg.Save()
		}
	}

	drawer := utils.NewDrawer(&cfg.DrawingType, &cfg.PositionX, &cfg.PositionY)
	searcher := utils.NewSearcher()

	ui.ShowUI(cfg, logger, drawer, searcher)
	hook.End()
}
