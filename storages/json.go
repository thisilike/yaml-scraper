package storages

import (
	"encoding/json"
	"errors"
	"os"
)

type JsonItem struct {
	Value string
	Name  string
	Type  string
}

type JsonStorage struct {
	Path            string
	Name            string
	Indent          string
	Prefix          string
	Template        []JsonItem
	StoreIncomplete bool
	JsonData        map[string]interface{}
	File            *os.File
}

func (s *JsonStorage) Open() error {
	var err error
	s.File, err = os.Create(s.Path)
	if err != nil {
		log.Errorf("failed to open storage file '%s' with error '%s'", s.Path, err)
		return errors.New("failed to open storage file")
	}
	s.JsonData = map[string]interface{}{}
	return nil
}

func (s *JsonStorage) Close() error {
	defer s.File.Close()
	j, err := json.MarshalIndent(s.JsonData, s.Prefix, s.Indent)
	if err != nil {
		log.Errorf("failed to marshal json '%s' with error", s.Name, err)
		return errors.New("failed to marshal json")
	}
	_, err = s.File.Write(j)
	if err != nil {
		log.Errorf("failed to write json byte's to file '%s' with error '%s'", s.Path, err)
		return errors.New("failed to write json byte's to file")
	}
	return nil
}

func (s *JsonStorage) Store(data map[string]interface{}) error {
	jsonInstanceData := map[string]interface{}{}
	for _, item := range s.Template {
		jsonInstanceData[item.Name] = data[item.Value]
	}
	s.JsonData[data["sourceUrl"].(string)] = jsonInstanceData
	return nil
}

func (s *JsonStorage) GetName() string {
	return s.Name
}

func (s *JsonStorage) Validate() error {
	return nil
}

func newJsonStorage(cnf map[string]interface{}) Storage {
	storage := JsonStorage{}

	// set path
	storage.Path = cnf["path"].(string)
	// set name
	if name, ok := cnf["name"]; ok {
		storage.Name = name.(string)
	} else {
		storage.Name = cnf["type"].(string)
	}
	// set indent
	if indent, ok := cnf["indent"]; ok {
		storage.Indent = indent.(string)
	} else {
		storage.Prefix = ""
	}
	// set prefix
	if prefix, ok := cnf["prefix"]; ok {
		storage.Prefix = prefix.(string)
	}
	// set template
	storage.StoreIncomplete = cnf["incomplete"].(bool)
	for _, item := range cnf["template"].([]interface{}) {
		it := item.(map[string]interface{})
		if _, exists := it["name"]; !exists {
			it["name"] = it["value"].(string)
		}
		storage.Template = append(storage.Template, JsonItem{
			Value: it["value"].(string),
			Name:  it["name"].(string),
		})
	}

	return &storage
}
