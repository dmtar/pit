package lib

import (
	"fmt"
	"regexp"
	"strings"
)

const EMAIL_REGEX string = `^([^@\s]+)@((?:[-a-z0-9]+\.)+[a-z]{2,})$`

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
		if val, ok := val.([]interface{}); ok {
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
		if ok := p.exists(p, key); !ok {
			return fmt.Errorf("The parameter %s is required!", key)
		}
	}

	return nil
}

func (p Params) ShouldBeEmail(key string) error {

	if ok, _ := regexp.MatchString(EMAIL_REGEX, p.Get(key)); !ok {
		return fmt.Errorf("Please provide valid email address!")
	}

	return nil
}

func (p Params) exists(input map[string]interface{}, key string) (ok bool) {
	if input == nil {
		return false
	}

	if index := strings.Index(key, "."); index != -1 {
		pair := strings.SplitN(key, ".", 2)
		params, asserted := input[pair[0]].(map[string]interface{})
		if !asserted {
			return false
		}
		ok = p.exists(params, pair[1])
	} else {
		_, ok = input[key]
	}

	return
}
