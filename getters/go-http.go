package getters

import (
	"errors"
	"io"
	"net/http"
)

type GoHttp struct {
	Name   string
	Retrys int
	Delay  int
}

func (g GoHttp) Open() error {
	return nil
}

func (g GoHttp) Close() error {
	return nil
}

func (g GoHttp) Get(url string) ([]byte, error) {
	c := 0
	for c < g.Retrys {
		resp, err := http.Get(url)
		if err != nil {
			if c < g.Retrys-1 {
				log.Warn(err)
				log.Warnf("'%s' failed getting data: '%s'", g.Name, url)
				log.Warn("retrying...")
				c++
				continue
			}
			log.Error(err)
			log.Errorf("'%s' failed getting data: '%s'", g.Name, url)
			return nil, errors.New("failed getting data")
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if c < g.Retrys-1 {
				log.Warn(err)
				log.Warnf("'%s' failed reading data: '%s'", g.Name, url)
				log.Warn("retrying...")
				c++
				continue
			}
			log.Error(err)
			log.Errorf("'%s' failed reading data: '%s'", g.Name, url)
			return nil, errors.New("failed reading data")
		}
		return body, nil
	}
	return nil, errors.New("should never be here")
}

func (g GoHttp) GetName() string {
	return g.Name
}

func (g GoHttp) Validate() error {
	if g.Name == "" {
		log.Errorf("empty getter name '%s' in go-http", "go-http")
		return errors.New("empty getter name")
	}
	if g.Delay < 0 {
		log.Errorf("no negative delay allowed '%s' in go-http", g.Delay)
		return errors.New("no negative delay allowed")
	}
	if g.Retrys < 0 {
		log.Errorf("no negative retry's allowed '%s' in go-http", g.Retrys)
		return errors.New("no negative retry's allowed")
	}

	return nil
}

func NewGoGetter(cnf map[string]interface{}) Getter {
	getter := GoHttp{}

	getter.Name, _ = cnf["name"].(string)
	getter.Delay, _ = cnf["delay"].(int)
	getter.Retrys, _ = cnf["retrys"].(int)

	return getter
}
