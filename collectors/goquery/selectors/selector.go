package selectors

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var SelectorRegister map[string]func(cnf map[string]interface{}) Selector

type Selector interface {
	Validate() error
	Select(*goquery.Selection) *goquery.Selection
}

func init() {
	log = config.Logger
	SelectorRegister = make(map[string]func(cnf map[string]interface{}) Selector)
	RegisterExtractor("css", newCssSelectorFunc)
}

func RegisterExtractor(name string, newSelectorFunc func(cnf map[string]interface{}) Selector) error {
	if _, ok := SelectorRegister[name]; ok {
		log.Errorf("selector with name '%s' already exist", name)
		return errors.New("selector does already exist")
	}
	SelectorRegister[name] = newSelectorFunc
	return nil
}

func NewSelector(cnf map[string]interface{}) Selector {
	return SelectorRegister[cnf["type"].(string)](cnf)
}
