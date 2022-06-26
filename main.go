package main

import (
	"fmt"

	"github.com/alitto/pond"
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
	"github.com/thisilike/ymls/scraper"
)

var log *logrus.Logger
var WorkerPool *pond.WorkerPool

func main() {
	scrapers, err := scraper.LoadConfig("ymls.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	// pool := pond.New(len(scrapers), 0)
	// for _, scraper := range scrapers {
	// 	pool.Submit(func() {
	// 		scraper.Start()
	// 	})
	// }
	// pool.StopAndWait()
	for _, scraper := range scrapers {
		scraper.Start()
	}
}

func init() {
	log = config.Logger
}
