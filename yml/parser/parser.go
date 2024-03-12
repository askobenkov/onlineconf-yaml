package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type OnlineConfItem struct {
	Key   string
	Value string
	Type  string
}

func GetParentNodeKeys(config map[string]OnlineConfItem) []string {
	nodes := map[string]bool{}
	for k := range config {

		list := strings.Split(k, "/")

		nodeKey := ""
		for _, k := range list[:len(list)-1] {
			if nodeKey == "" {

				nodeKey = k
			} else {

				nodeKey = nodeKey + "/" + k
			}
			if ok := nodes[nodeKey]; !ok {
				nodes[nodeKey] = true
			}
		}
	}
	nodeKeys := []string{}
	for k := range nodes {
		nodeKeys = append(nodeKeys, k)
	}

	sort.Strings(nodeKeys)
	return nodeKeys
}

func GetNodeKeysForDelete(config map[string]OnlineConfItem) []string {
	nodes := map[string]bool{}
	for k := range config {

		list := strings.Split(k, "/")

		nodeKey := ""
		for _, k := range list {
			if nodeKey == "" {

				nodeKey = k
			} else {

				nodeKey = nodeKey + "/" + k
			}
			if ok := nodes[nodeKey]; !ok {
				nodes[nodeKey] = true
			}
		}
	}
	nodeKeys := []string{}
	for k := range nodes {
		nodeKeys = append(nodeKeys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(nodeKeys)))
	return nodeKeys
}

func WalkByYML(obj reflect.Value, prefix string) map[string]OnlineConfItem {
	o := make(map[string]OnlineConfItem)
	switch obj.Kind() {
	case reflect.Ptr:
		res := WalkByYML(obj.Elem(), "")
		o = mergeMaps(o, res)
	case reflect.Interface:
		res := WalkByYML(obj.Elem(), prefix)
		o = mergeMaps(o, res)
	case reflect.Map:
		if prefix != "" {
			childrenKeys := []string{}
			for _, key := range obj.MapKeys() {
				p := key.Elem().String()
				childrenKeys = append(childrenKeys, p)
			}
			jsonBytes, err := json.Marshal(childrenKeys)
			if err != nil {
				log.Printf("Can't marshal node keys for the '%s': %s", prefix, err.Error())
				break
			}
			o[prefix] = OnlineConfItem{
				Key:   prefix,
				Value: string(jsonBytes),
				Type:  "application/json",
			}
		}
		for _, key := range obj.MapKeys() {
			p := key.Elem().String()
			if prefix != "" {
				//				keyValue := ""
				//				switch key.Elem() {
				//					case reflect.Int
				//				}
				p = fmt.Sprintf("%s/%v", prefix, key.Elem())
				//				log.Printf("================> %+v", key.Elem())
				//				log.Printf("================> %+v", key.Elem())
			}
			res := WalkByYML(obj.MapIndex(key), p)
			//			o[prefix] = res
			o = mergeMaps(o, res)
		}
	case reflect.Slice:
		////		o[prefix] = obj.Interface()
		//list := []interface{}{}
		list := []string{}
		for i := 0; i < obj.Len(); i++ {
			sliceObj := obj.Index(i)
			toMarshalIf := sliceObj.Interface()
			switch v := toMarshalIf.(type) {
			case string, float64, int, bool:
				list = append(list, fmt.Sprintf("- %v", v))
			default:
				toMarshalIf = []interface{}{v}
				yamlOption, err := yaml.Marshal(toMarshalIf)
				if err != nil {
					log.Fatal(fmt.Errorf("can't marshal %+v to yaml... %+v", sliceObj.Interface(), err))
				}
				list = append(list, fmt.Sprintf("%v", string(yamlOption)))
			}
		}
		if len(list) > 0 {
			o[prefix] = OnlineConfItem{
				Key:   prefix,
				Value: strings.Join(list, "\n"),
				Type:  "application/x-yaml",
			}
		}
	case reflect.String, reflect.Float64, reflect.Int, reflect.Bool:
		o[prefix] = OnlineConfItem{
			Key:   prefix,
			Value: fmt.Sprintf("%v", obj.Interface()),
			Type:  "text/plain",
		}
	}
	return o
}

func GetYMLConfig(filepath string) (interface{}, error) {
	var data interface{}

	readFile, err := os.Open(filepath)
	if err != nil {
		return data, err
	}
	defer readFile.Close()
	reader := bufio.NewReader(readFile)
	yamlProcessed := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return data, err
		}
		yamlProcessed += line
	}

	err = yaml.Unmarshal([]byte(yamlProcessed), &data)
	if err != nil {
		return data, err
	}

	return data, err
}

func mergeMaps(a, b map[string]OnlineConfItem) map[string]OnlineConfItem {
	for k, v := range b {
		a[k] = v
	}
	return a
}
