package reporter

import (
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger

type Reporter interface {
	Open() (err error)
	Close() (err error)
	Report(map[string]interface{})
	GetName() string
	Validate() error
}

func init() {
	log = config.Logger
}
