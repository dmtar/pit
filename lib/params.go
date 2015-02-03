package lib

import "fmt"

const EMAIL_REGEX string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

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

func (p Params) ShouldBeEmail(str string) error {
	rxEmail = regexp.MustCompile(EMAIL_REGEX)
	return rxEmail.MatchString(str)
}
