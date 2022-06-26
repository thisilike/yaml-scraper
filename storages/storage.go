package storages

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var StorageRegister map[string]func(map[string]interface{}) Storage

type Storage interface {
	Open() (err error)
	Close() (err error)
	Store(map[string]interface{}) (err error)
	GetName() string
	Validate() error
}

func init() {
	log = config.Logger
	StorageRegister = make(map[string]func(map[string]interface{}) Storage)
	RegisterStorage("json", newJsonStorage)
}

func RegisterStorage(name string, newStorage func(map[string]interface{}) Storage) error {
	if _, ok := StorageRegister[name]; ok {
		log.Errorf("storage with name '%s' already exist", name)
		return errors.New("storage does already exist")
	}
	StorageRegister[name] = newStorage
	return nil
}

func NewStorage(cnf map[string]interface{}) Storage {
	return StorageRegister[cnf["type"].(string)](cnf)
}
