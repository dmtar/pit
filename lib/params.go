package lib

import "fmt"

type Params map[string]interface{}

func (p Params) Get(key string) string {
	if val, ok := p[key]; ok {
		return fmt.Sprintf("%v", val)
	}

	return ""
}

func (p Params) Required(keys ...string) error {
	for _, key := range keys {
		if _, ok := p[key]; !ok {
			return fmt.Errorf("The parameter %s is required!", key)
		}
	}
	return nil
}
