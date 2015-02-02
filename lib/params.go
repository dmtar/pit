package lib

import "fmt"

type Params map[string]interface{}

func (p Params) Get(key string) string {
	if val, ok := p[key]; ok {
		return fmt.Sprintf("%v", val)
	}

	return ""
}

func (p Params) GetP(key string) Params {
	if val, ok := p[key]; ok {
		return Params(val.(map[string]interface{}))
	}

	return Params{}
}

func (p Params) GetA(key string) []interface{} {

	if val, ok := p[key]; ok {
		return val.([]interface{})
	}

	return nil
}

func (p Params) GetAString(key string) []string {
	result := make([]string, 0)
	if val, ok := p[key]; ok {
		val, ok := val.([]interface{})
		if ok {
			for _, v := range val {
				if vs, ok := v.(string); ok {
					result = append(result, vs)
				}
			}
		}
		return result
	}

	return nil
}

func (p Params) Required(keys ...string) error {
	for _, key := range keys {
		if _, ok := p[key]; !ok {
			return fmt.Errorf("The parameter %s is required!", key)
		}
	}
	return nil
}
