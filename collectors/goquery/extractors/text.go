package extractors

import "github.com/PuerkitoBio/goquery"

type TextExtractor struct {
}

func (ex *TextExtractor) Extract(sel *goquery.Selection) interface{} {
	return sel.Text()
}

func (ex *TextExtractor) Validate() error {
	return nil
}

func newTextExtractorFunc(cnf map[string]interface{}) Extractor {
	extractor := TextExtractor{}
	return &extractor
}
