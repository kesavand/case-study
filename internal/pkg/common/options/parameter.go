package options

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type UserInfo struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	CardNo    string `json:"cardno"`
	Timestamp string `json:"timestamp"`
}

//ParameterBag construction for holding configuration options
type ParameterBag struct {
	parameters *sync.Map
}

//New creates empty bag
func New() *ParameterBag {
	bag := &ParameterBag{
		parameters: &sync.Map{},
	}

	return bag

}

//Read reads interface value, if not found, will read from envs, if not found there will return defaultVal
func (p *ParameterBag) Read(name string, defaultVal interface{}) (interface{}, bool) {
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

//ReadDuration reads int value and converts it to duration identified by the unit, if not set, will return defaultVal
func (p *ParameterBag) ReadDuration(name string, unit time.Duration, defaultVal uint) time.Duration {
	val := p.ReadUint(name, defaultVal)
	return unit * time.Duration(val)
}

//ReadUint same as Read but returns a uint
func (p *ParameterBag) ReadUint(name string, defaultVal uint) uint {
	valI, found := p.Read(name, defaultVal)
	if !found {
		return defaultVal
	}

	val, ok := valI.(uint)
	if !ok {
		uintVal, err := strconv.ParseUint(fmt.Sprint(valI), 10, 32)
		if err == nil {
			return uint(uintVal)
		}
		return defaultVal
	}

	return val
}
