package scraper

import (
	"errors"
	"time"

	"github.com/alitto/pond"
	"github.com/thisilike/ymls/getters"
)

type Scraper struct {
	Name       string
	Urls       []string
	Getters    []getters.Getter
	Processors []Processor
	MaxWorkers int
	Delay      int
}

type void struct{}

func (scraper *Scraper) Validate() error {
	// name
	if _, exists := ScraperRegister[scraper.Name]; exists {
		log.Errorf("scraper name duplicate '%s'", scraper.Name)
		return errors.New("scraper name duplicate")
	}
	ScraperRegister[scraper.Name] = void{}

	// urls not validated for now.
	// Getters
	for _, getter := range scraper.Getters {
		err := getter.Validate()
		if err != nil {
			log.Errorf("invalid getters '%s'", err)
			return errors.New("invalid getters")
		}
	}
	// Processors
	for _, processor := range scraper.Processors {
		err := processor.Validate()
		if err != nil {
			log.Errorf("invalid processors '%s'", err)
			return errors.New("invalid processors")
		}
	}

	// MaxWorkers
	if scraper.MaxWorkers < 0 {
		log.Errorf("maxWorkers must not be negative '%s'", scraper.MaxWorkers)
		return errors.New("maxWorkers must not be negative")
	} else if scraper.MaxWorkers == 0 {
		log.Warn("maxWorkers not set. using defaults (1)")
	}

	// Delay
	if scraper.Delay < 0 {
		log.Errorf("delay must not be negative '%s'", scraper.Delay)
		return errors.New("delay must not be negative")
	}

	return nil
}

func (scraper *Scraper) Start() error {
	// open storage interfaces
	for _, processor := range scraper.Processors {
		for _, store := range processor.Storages {
			err := store.Open()
			if err != nil {
				log.Errorf("failed to open storage '%s'", store.GetName())
				return errors.New("failed to open storage")
			}
			defer store.Close()
		}
	}

	// star scraping
	pool := pond.New(int(scraper.MaxWorkers), len(scraper.Urls))
	f := func(url string) func() {
		return func() {
			scraper.Scrape(url)
		}
	}
	func() {
		for _, url := range scraper.Urls {
			pool.Submit(f(url))
			time.Sleep(time.Duration(scraper.Delay) * time.Millisecond)
		}
	}()
	pool.StopAndWait()

	log.Infof(
		"Scraper: '%s' finished!", scraper.Name)
	return nil
}

func (scraper *Scraper) Scrape(url string) error {
	// get data
	data, err := scraper.Get(url)
	if err != nil {
		return errors.New("failed to get data from url")
	}

	// process data
	err = scraper.Process(data, url)
	if err != nil {
		return errors.New("failed processing data")
	}
	return nil
}

func (scraper *Scraper) Get(url string) ([]byte, error) {
	log.Infof("getting data from url '%s'", url)
	var (
		data []byte
		err  error
	)
	for i, getter := range scraper.Getters {
		data, err = getter.Get(url)
		if err != nil {
			log.Warnf("failed to get with getter '%s'", getter.GetName())
			if i == len(scraper.Getters)-1 {
				log.Errorf("failed to get data for url '%s'", url)
				return data, errors.New("failed to get data")
			}
			continue
		}
		break
	}
	return data, nil
}

func (scraper *Scraper) Process(data []byte, url string) error {
	processPool := pond.New(len(scraper.Processors), 1)
	spawnProcess := func(processor Processor) func() {
		return func() {
			err := processor.Process(data, url)
			if err != nil {
				log.Errorf("processing incomplete. scraper:'%s'", scraper.Name)
			}
		}
	}
	for _, processor := range scraper.Processors {
		processPool.Submit(spawnProcess(processor))
	}
	processPool.StopAndWait()
	return nil
}

func NewScraper(cnf map[string]interface{}) Scraper {
	scraper := Scraper{}
	scraper.Name = cnf["name"].(string)
	scraper.MaxWorkers = cnf["maxWorkers"].(int)
	scraper.Delay = cnf["delay"].(int)
	urls := cnf["urls"].([]interface{})
	for _, v := range urls {
		scraper.Urls = append(scraper.Urls, v.(string))
	}
	scraper.Getters = NewGetters(cnf["getters"].([]interface{}))
	scraper.Processors = NewProcessors(cnf["processors"].([]interface{}))
	return scraper
}

func NewGetters(cnf []interface{}) []getters.Getter {
	var gettersList []getters.Getter
	for _, c := range cnf {
		gettersList = append(gettersList, getters.NewGetter(c.(map[string]interface{})))
	}
	return gettersList
}

func NewProcessors(cnf []interface{}) []Processor {
	var processorList []Processor
	for _, c := range cnf {
		processorList = append(processorList, NewProcessor(c.(map[string]interface{})))
	}
	return processorList
}
