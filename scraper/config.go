package scraper

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
	"gopkg.in/yaml.v3"
)

func ValidateConfigSchema(config interface{}) error {
	compiler := jsonschema.NewCompiler()
	rootDir, err := os.Getwd()
	rootDir = rootDir + "/schema"
	if err != nil {
		return err
	}
	err = filepath.Walk(rootDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return err
		}

		schemaBytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Errorf("failed to read schema: 'ymls.json'", err.Error())
			return err
		}

		if err := compiler.AddResource(rootDir+path, strings.NewReader(string(schemaBytes))); err != nil {
			log.Errorf("failed to add resource to schema compiler: '%s'", err.Error())
			return errors.New("failed to add schema compiler resource")
		}

		return err
	})
	if err != nil {
		return err
	}

	schema, err := compiler.Compile("./schema/ymls.json")
	if err != nil {
		log.Errorf("failed to compile schema: '%s'", err.Error())
		return errors.New("failed to compile schema")
	}

	if err := schema.ValidateInterface(config); err != nil {
		log.Errorf("invalid config: '%s'", err.Error())
		return errors.New("invalid config")
	}

	log.Debug("Config Valid!")
	return nil
}

func LoadConfig(path string) ([]Scraper, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	scrapers, err := deSerializeConfig(data)
	if err != nil {
		log.Errorf("failed to parse config: '%s'", err)
		return nil, errors.New("failed loading config")
	}
	return scrapers, err
}

func deSerializeConfig(data []byte) ([]Scraper, error) {
	var s map[interface{}]interface{}
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		log.Errorf("failed parsing config: '%s'", err)
		return nil, errors.New("failed parsing config")
	}

	m, err := toStringKeys(s)
	if err != nil {
		return nil, err
	}

	err = ValidateConfigSchema(m)
	if err != nil {
		return nil, err
	}

	d := m.(map[string]interface{})
	l := d["scrapers"].([]interface{})
	var scrapers []Scraper
	for _, scraperCnf := range l {
		scrapers = append(scrapers, NewScraper(scraperCnf.(map[string]interface{})))
	}

	for _, scraper := range scrapers {
		if err := scraper.Validate(); err != nil {
			log.Errorf("invalid config: '%s'")
			return nil, err
		}
	}

	return scrapers, nil
}

func toStringKeys(val interface{}) (interface{}, error) {
	var err error
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				log.Errorf("non string key '%s' of type '%T' as as key are not supported", k, k)
				return nil, errors.New("found non-string key")
			}
			m[k], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []interface{}:
		var l = make([]interface{}, len(val))
		for i, v := range val {
			l[i], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return val, nil
	}
}
