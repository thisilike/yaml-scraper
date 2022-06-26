package getters

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var GetterRegister map[string]func(map[string]interface{}) Getter

type Getter interface {
	Open() (err error)
	Close() (err error)
	Get(url string) (data []byte, err error)
	GetName() string
	Validate() error
}

type GetterReport struct {
}

func init() {
	log = config.Logger
	GetterRegister = make(map[string]func(map[string]interface{}) Getter)
	RegisterGetter("go-http", NewGoGetter)
}

func RegisterGetter(name string, newGetterF func(map[string]interface{}) Getter) error {
	if _, exists := GetterRegister[name]; exists {
		log.Errorf("getter name duplicate '%s'", name)
		return errors.New("getter name duplicate")
	}
	GetterRegister[name] = newGetterF
	return nil
}

func NewGetter(cnf map[string]interface{}) Getter {
	return GetterRegister[cnf["type"].(string)](cnf)
}
