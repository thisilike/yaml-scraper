package selectors

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

type CssSelector struct {
	cssSelectorString string
}

func (selector *CssSelector) Validate() error {
	if _, err := cascadia.Parse(selector.cssSelectorString); err != nil {
		log.Errorf("'%s' is an invalid css selector. '%s'", selector.cssSelectorString, err)
		return errors.New("invalid css selector")
	}
	return nil
}

func (s *CssSelector) Select(sel *goquery.Selection) *goquery.Selection {
	return sel.Find(s.cssSelectorString)
}

func newCssSelectorFunc(cnf map[string]interface{}) Selector {
	selector := CssSelector{}
	selector.cssSelectorString = cnf["selectorString"].(string)
	return &selector
}
