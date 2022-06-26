package scraper

import (
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var ScraperRegister map[string]struct{}

func init() {
	log = config.Logger
	ScraperRegister = make(map[string]struct{})
}
