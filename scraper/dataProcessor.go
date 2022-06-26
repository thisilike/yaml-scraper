package scraper

import (
	"errors"

	"github.com/thisilike/ymls/collectors"
	"github.com/thisilike/ymls/storages"
	"github.com/thisilike/ymls/transformers"
)

type Processor struct {
	Collectors   []collectors.Collector
	Storages     []storages.Storage
	Transformers []transformers.Transformer
	Data         map[string]interface{}
}

func (p *Processor) Validate() error {
	// validate collectors
	tmp := map[string]bool{}
	for _, collector := range p.Collectors {
		if tmp[collector.GetName()] {
			log.Errorf("no duplicate names allowed for collector(%s) in processor)", collector.GetName())
			return errors.New("no duplicate collector names allowed for processor")
		}
		err := collector.Validate()
		if err != nil {
			log.Errorf("invalid collector '%s'", collector.GetName())
			return errors.New("invalid collector")
		}
	}

	// validate storages
	for _, storage := range p.Storages {
		err := storage.Validate()
		if err != nil {
			log.Errorf("invalid storage '%s'", storage.GetName())
			return errors.New("invalid storage")
		}
	}

	// validate transformers
	for _, transformer := range p.Transformers {
		err := transformer.Validate()
		if err != nil {
			log.Errorf("invalid transformer")
			return errors.New("invalid transformer")
		}
	}

	return nil
}

func (p *Processor) Process(sourceData []byte, url string) error {
	p.Data = map[string]interface{}{}
	p.Data["sourceUrl"] = url
	err := p.Collect(sourceData)
	if err != nil {
		log.Errorf("collecting data incomplete '%s'", err)
	}
	err = p.Transform()
	if err != nil {
		log.Errorf("transforming data incomplete '%s'", err)
	}
	err = p.Store()
	if err != nil {
		log.Errorf("storing data incomplete '%s'", err)
	}
	return nil
}

func (p *Processor) Transform() error {
	for _, transformer := range p.Transformers {
		if _, ok := p.Data[transformer.Name]; !ok {
			log.Errorf("could not find variable with transformer name: '%s'", transformer.Name)
			continue
		}
		val := transformer.DoTransformations(p.Data[transformer.Name])
		p.Data[transformer.Name] = val
	}
	return nil
}

func (p *Processor) Collect(data []byte) error {
	var err error
	for _, collector := range p.Collectors {
		// open
		err = collector.Open()
		if err != nil {
			log.Errorf("failed to open collector '%s'", collector.GetName())
		}
		defer collector.Close()

		// collect
		var collData interface{}
		collData, err = collector.Collect(data)
		if err != nil {
			log.Errorf("failed to collect data with '%s'", collector.GetName())
			continue
		}
		p.Data[collector.GetName()] = collData
	}
	if err != nil {
		return errors.New("collecting data incomplete")
	}
	return nil
}

func (p *Processor) Store() error {
	var err error
	for _, store := range p.Storages {
		// store
		err = store.Store(p.Data)
		if err != nil {
			log.Errorf("failed to store in storage '%s'", store.GetName())
			continue
		}
	}
	if err != nil {
		return errors.New("storing is incomplete")
	}
	return nil
}

func NewProcessor(cnf map[string]interface{}) Processor {
	processor := Processor{}
	processor.Collectors = NewCollectors(cnf["collectors"].([]interface{}))
	if _, ok := cnf["transformers"].([]interface{}); ok {
		processor.Transformers = NewTransformers(cnf["transformers"].([]interface{}))
	}
	processor.Storages = NewStorages(cnf["storages"].([]interface{}))
	return processor
}

func NewCollectors(cnf []interface{}) []collectors.Collector {
	var collectorsList []collectors.Collector
	for _, c := range cnf {
		collectorsList = append(collectorsList, collectors.NewCollector(c.(map[string]interface{})))
	}
	return collectorsList
}

func NewTransformers(cnf []interface{}) []transformers.Transformer {
	var transformersList []transformers.Transformer
	for _, c := range cnf {
		transformersList = append(transformersList, transformers.NewTransformer(c.(map[string]interface{})))
	}
	return transformersList
}

func NewStorages(cnf []interface{}) []storages.Storage {
	var storagesList []storages.Storage
	for _, c := range cnf {
		storagesList = append(storagesList, storages.NewStorage(c.(map[string]interface{})))
	}
	return storagesList
}
