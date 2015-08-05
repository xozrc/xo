package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var print = fmt.Print

type StoreConfig struct {
	Store  string          `json:"store"`
	Config json.RawMessage `json:"config"`
}

func NewStoreConf(path string) (conf *StoreConfig, err error) {
	file, err1 := os.Open(path)
	if err1 != nil {
		err = err1
		return
	}
	defer file.Close()
	content, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		err = err2
		return
	}
	conf = &StoreConfig{}

	err = json.Unmarshal(content, conf)

	return
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

func NewServerConf(m json.RawMessage) (conf *ServerConfig, err error) {
	conf = &ServerConfig{}
	err = json.Unmarshal(m, conf)
	return
}
