package collectors

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var CollectorRegister map[string]func(map[string]interface{}) Collector

type Collector interface {
	Open() (err error)
	Close() (err error)
	Collect(source []byte) (data interface{}, err error)
	GetName() string
	Validate() error
}

func init() {
	log = config.Logger
	CollectorRegister = make(map[string]func(map[string]interface{}) Collector)
	RegisterCollector("go-query", newGoQueryCollector)
}

func RegisterCollector(name string, newCollectorF func(map[string]interface{}) Collector) error {
	if _, ok := CollectorRegister[name]; ok {
		log.Errorf("collector with name '%s' already exist", name)
		return errors.New("collector does already exist")
	}
	CollectorRegister[name] = newCollectorF
	return nil
}

func NewCollector(cnf map[string]interface{}) Collector {
	return CollectorRegister[cnf["type"].(string)](cnf)
}
