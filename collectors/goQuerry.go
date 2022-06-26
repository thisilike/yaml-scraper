package collectors

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/thisilike/ymls/collectors/goquery/extractors"
	"github.com/thisilike/ymls/collectors/goquery/selectors"
)

type GoQueryCollector struct {
	name      string
	selectors []selectors.Selector
	extractor extractors.Extractor
}

func (c *GoQueryCollector) Open() error {
	return nil
}

func (c *GoQueryCollector) Close() error {
	return nil
}

func (c *GoQueryCollector) Collect(source []byte) (interface{}, error) {
	// parse
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(source)))
	if err != nil {
		log.Error("failed parsing html in go-querry reader '%s'", err)
		return nil, errors.New("failed parsing html")
	}
	// select
	s := doc.Find("html")
	for _, selector := range c.selectors {
		s = selector.Select(s)
	}
	// extract
	return c.extractor.Extract(s), nil
}

func (c *GoQueryCollector) GetName() string {
	return c.name
}

func (c *GoQueryCollector) Validate() error {
	for _, sel := range c.selectors {
		err := sel.Validate()
		if err != nil {
			log.Errorf("invalid selector configuration in collector '%s'", c.name)
			return errors.New("invalid selector configuration in collector")
		}
	}
	err := c.extractor.Validate()
	if err != nil {
		log.Errorf("invalid extractor configuration in collector '%s'", c.name)
		return errors.New("invalid extractor configuration")
	}
	return nil
}

func newGoQueryCollector(cnf map[string]interface{}) Collector {
	coll := GoQueryCollector{}
	// set name
	coll.name = cnf["name"].(string)
	// set selectors
	selectorsList := cnf["selectors"].([]interface{})
	for _, sel := range selectorsList {
		coll.selectors = append(coll.selectors, selectors.NewSelector(sel.(map[string]interface{})))
	}
	// set extractor
	coll.extractor = extractors.NewExtractor(cnf["extractor"].(map[string]interface{}))

	return &coll
}
