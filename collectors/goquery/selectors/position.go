package selectors

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

type PositionSelector struct {
	Index int
}

func (s *PositionSelector) Validate() error {
	if s.Index < 0 {
		log.Errorf("position index mus be positive or 0 not '%s'", s.Index)
		return errors.New("invalid position index")
	}
	return nil
}

func (s *PositionSelector) Select(sel *goquery.Selection) *goquery.Selection {
	// TODO: add last element
	selection := &goquery.Selection{}
	sel.Map(func(i int, sele *goquery.Selection) string {
		if i == s.Index {
			selection = sele
		}
		return ""
	})
	return selection
}

func newPositionSelector(cnf map[string]interface{}) Selector {
	selector := PositionSelector{}
	selector.Index = cnf["index"].(int)
	return &selector
}
