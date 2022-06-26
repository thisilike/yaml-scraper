package transformers

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var TransformationRegister map[string]func(map[string]interface{}) Transformation

func init() {
	log = config.Logger
	TransformationRegister = make(map[string]func(map[string]interface{}) Transformation)
}

func RegisterTransformation(name string, newTransformationFunc func(map[string]interface{}) Transformation) error {
	if _, exists := TransformationRegister[name]; exists {
		log.Errorf("Transformation with the name '%s' does already exist", name)
		return errors.New("transformation name duplication")
	}
	TransformationRegister[name] = newTransformationFunc
	return nil
}

func InStringList(v string, list []string) bool {
	for _, val := range list {
		if val == v {
			return true
		}
	}
	return false
}
