package cdb

import (
	"encoding/json"
	"math"
	"regexp"
	"strings"

	"github.com/colinmarc/cdb"
	"gopkg.in/yaml.v2"
)

var floatRE = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][-+]?\d+)?$`) // RFC8259

// WriteItem cdb struct
type WriteItem struct {
	json  bool
	Path  string
	Value string
	Tp    string
}

type yamlValue struct {
	value interface{}
	str   string
}

func (v *yamlValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// handle map[interface{}]interface{} separately to prevent unwanted recursion
	{
		var value map[string]*yamlValue
		err := unmarshal(&value)
		if err == nil {
			v.value = value
			return nil
		}
	}

	// handle []interface{} separately to prevent unwanted recursion
	{
		var value []*yamlValue
		err := unmarshal(&value)
		if err == nil {
			v.value = value
			return nil
		}
	}

	var value interface{}
	err := unmarshal(&value)
	if err != nil {
		return err
	}

	switch typedVal := value.(type) {
	case bool: // strict YAML 1.2
		var str string
		err := unmarshal(&str)
		if err != nil {
			return err
		}
		switch str {
		case "true", "false":
			v.value = value
		default:
			v.value = str
		}
		return nil
	case float64:
		if math.IsNaN(typedVal) || math.IsInf(typedVal, 0) {
			v.value = nil
		} else {
			v.value = value
			unmarshal(&v.str)
			if !floatRE.MatchString(v.str) {
				v.str = ""
			}
		}
		return nil
	default:
		v.value = value
		return nil
	}
}

func (v *yamlValue) MarshalJSON() ([]byte, error) {
	if v.str != "" {
		return []byte(v.str), nil
	}
	return json.Marshal(v.value)
}

// YAMLToJSON - convert yaml to json
func YAMLToJSON(y []byte) ([]byte, error) {
	var data *yamlValue
	err := yaml.Unmarshal(y, &data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

// Write item to filepath
func Write(filepath string, params []WriteItem) error {
	w, err := cdb.Create(filepath)
	if err != nil {
		return err
	}

	for _, param := range params {
		switch param.Tp {
		case "application/x-yaml":
			res, err := YAMLToJSON([]byte(param.Value))
			if err != nil {
				return err
			}
			param.Value = string(res)
			param.json = true
		}

		p := strings.Split(param.Path, "/")
		param.Path = strings.Join(p, ".")
		var t string
		if param.json {
			t = "j"
		} else {
			t = "s"
		}
		err = w.Put([]byte(param.Path), []byte(t+param.Value))
		if err != nil {
			w.Close()
			return err
		}
	}

	return w.Close()
}
