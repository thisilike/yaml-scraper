package extractors

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var ExtractorRegister map[string]func(cnf map[string]interface{}) Extractor

type Extractor interface {
	Validate() error
	Extract(*goquery.Selection) interface{}
}

func init() {
	log = config.Logger
	ExtractorRegister = make(map[string]func(cnf map[string]interface{}) Extractor)
	RegisterExtractor("text", newTextExtractorFunc)

}

func RegisterExtractor(name string, newExtractorFunc func(cnf map[string]interface{}) Extractor) error {
	if _, ok := ExtractorRegister[name]; ok {
		log.Errorf("extractor with name '%s' already exist", name)
		return errors.New("extractor does already exist")
	}
	ExtractorRegister[name] = newExtractorFunc
	return nil
}

func NewExtractor(cnf map[string]interface{}) Extractor {
	return ExtractorRegister[cnf["type"].(string)](cnf)
}
