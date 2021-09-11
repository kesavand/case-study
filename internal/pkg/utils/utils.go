package utils

import (
	"context"
	"os"
	"sync"
)

type UserInfo struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	CardNo    string `json:"cardno"`
	Timestamp string `json:"timestamp"`
}

//WaitForExit waits on an exit channel until error or an interrupt is received
func Wait(exitChannel chan error, cancel context.CancelFunc) error {
	err := <-exitChannel
	// cancel the context
	cancel()
	return err
}

//Parameter construction for holding configuration options
type Parameter struct {
	parameters *sync.Map
}

//New creates empty bag
func New() *Parameter {
	bag := &Parameter{
		parameters: &sync.Map{},
	}

	return bag

}

//Read reads from env variable s, returns default if not present in env
func (p *Parameter) Read(name string, defaultVal interface{}) (interface{}, bool) {
	valI, okI := p.parameters.Load(name)
	if okI {
		return valI, true
	}
	envVal, found := os.LookupEnv(name)
	if found {
		return envVal, true
	}
	return defaultVal, false
}

func (p *Parameter) ReadString(name string, defaultVal string) (string, bool) {
	val, found := p.Read(name, defaultVal)

	if !found {
		return defaultVal, false
	}

	return val.(string), true
}
