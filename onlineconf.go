package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"log"

	"gopkg.in/yaml.v2"
)

/*
  go run onlineconf.go -onlineConfUrl https://onlineconf.local/config -exportConfigFilepath ./Revise1.yml -headersFilepath ./headers.txt -mainNodeName revise -showParsedConfig
*/
func main() {

	onlineConfUrl := flag.String("onlineConfUrl", "https://onlineconf.local/config", "OnlineConf URL name")
	configFilepath := flag.String("exportConfigFilepath", "", "export config filepath")
	headersFilepath := flag.String("headersFilepath", "", "headers filepath")
	mainNodeName := flag.String("mainNodeName", "", "OnlineConf main node name")
	showParsedConfig := flag.Bool("showParsedConfig", false, "Show parsed config")
	exportParsedConfig := flag.Bool("exportParsedConfig", false, "Export parsed config to OnlineConf")
	deleteParsedConfig := flag.Bool("deleteParsedConfig", false, "Delete config in OnlineConf")
	skipAlreadyExist := flag.Bool("skipAlreadyExist", false, "Skip already exist error")
	skipCreateNode := flag.Bool("skipCreateNode", false, "Skip create node")

	flag.Parse()

	if *configFilepath == "" {
		log.Fatal(fmt.Errorf("export filepath config is empty"))
	}

	client, err := NewOnlineConfClient(
		fmt.Sprintf("%s/%s", *onlineConfUrl, *mainNodeName),
		*headersFilepath,
	)
	if err != nil {
		log.Fatal(err)
	}

	data, err := GetYMLConfig(*configFilepath)
	if err != nil {
		log.Fatal(err)
	}

	obj := reflect.ValueOf(&data)
	src := WalkByYML(obj, "")

	if *showParsedConfig {
		for k, v := range src {
			fmt.Printf("%-50s (%-30s) : %v\n", k, v.Type, v.Value)
		}
	}

	nodeKeys := getParentNodeKeys(src)
	log.Printf("create =============> %+v\n", nodeKeys)

	if *exportParsedConfig {
		if !*skipCreateNode {
			log.Printf("client.CreateEmptyNode..")

			for _, key := range nodeKeys {
				err := client.CreateEmptyNode(key, *skipAlreadyExist)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		for k, v := range src {

			log.Printf("%s (%s) -> %v\n", k, v.Type, v.Value)
			err := client.CreateNode(v)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if *deleteParsedConfig {
		nodeKeys = getNodeKeysForDelete(src)
		log.Printf("delete =============> %+v\n", nodeKeys)

		for _, key := range nodeKeys {
			err := client.DeleteNode(key)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return

}

type OnlineConfItem struct {
	Key   string
	Value string
	Type  string
}

func getParentNodeKeys(config map[string]OnlineConfItem) []string {
	nodes := map[string]bool{}
	for k, _ := range config {

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
	for k, _ := range nodes {
		nodeKeys = append(nodeKeys, k)
	}

	sort.Strings(nodeKeys)
	return nodeKeys
}

func getNodeKeysForDelete(config map[string]OnlineConfItem) []string {
	nodes := map[string]bool{}
	for k, _ := range config {

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
	for k, _ := range nodes {
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
		for _, key := range obj.MapKeys() {
			p := key.Elem().String()
			if prefix != "" {
				//				keyValue := ""
				//				switch key.Elem() {
				//					case reflect.Int
				//				}
				p = fmt.Sprintf("%s/%v", prefix, key.Elem())
				log.Printf("================> %+v", key.Elem())
				log.Printf("================> %+v", key.Elem())
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
			switch sliceObj.Kind() {
			case reflect.String, reflect.Float64, reflect.Int, reflect.Bool, reflect.Interface:
				list = append(list, fmt.Sprintf("%v", sliceObj.Interface()))
			default:
				log.Fatal(fmt.Errorf("unsuported type: %+v, %+v", sliceObj.Kind(), sliceObj))
			}
			//res := WalkByYML(obj.Index(i), fmt.Sprintf("%s.%d", prefix, i))
			//o = mergeMaps(o, res)
		}
		if len(list) > 0 {
			o[prefix] = OnlineConfItem{
				Key:   prefix,
				Value: strings.Join(list, ","),
				Type:  "application/x-list",
			}
		}
	case reflect.String:
		o[prefix] = OnlineConfItem{
			Key:   prefix,
			Value: fmt.Sprintf("%v", obj.Interface()),
			Type:  "text/plain",
		}
	case reflect.Float64:
		o[prefix] = OnlineConfItem{
			Key:   prefix,
			Value: fmt.Sprintf("%v", obj.Interface()),
			Type:  "text/plain",
		}
	case reflect.Int:
		o[prefix] = OnlineConfItem{
			Key:   prefix,
			Value: fmt.Sprintf("%v", obj.Interface()),
			Type:  "text/plain",
		}
	case reflect.Bool:
		o[prefix] = OnlineConfItem{
			Key:   prefix,
			Value: fmt.Sprintf("%v", obj.Interface()),
			Type:  "text/plain",
		}
	}
	return o
}

func mergeMaps(a, b map[string]OnlineConfItem) map[string]OnlineConfItem {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func GetYMLConfig(filepath string) (interface{}, error) {
	var data interface{}

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return data, err
	}

	return data, err
}

type OnlineConfClient struct {
	host    string
	headers map[string]string
}

type OnlineConfResponse struct {
	Error   string
	Version int
	Message string
}

func NewOnlineConfClient(
	host string,
	filepathHeader string,
) (*OnlineConfClient, error) {

	client := &OnlineConfClient{
		headers: map[string]string{},
		host:    host,
	}

	if filepathHeader != "" {
		headers, err := client.GetHeaders(filepathHeader)
		if err != nil {
			return nil, err
		}
		client.headers = headers
	}

	return client, nil
}

func (client *OnlineConfClient) GetHeaders(filepath string) (map[string]string, error) {

	headers := map[string]string{}
	if filepath == "" {
		return headers, nil
	}

	file, err := os.Open(filepath)
	if err != nil {
		return headers, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		list := strings.Split(line, ":")
		headers[strings.TrimSpace(list[0])] = strings.TrimSpace(list[1])
	}
	if err := scanner.Err(); err != nil {
		return headers, err
	}
	return headers, nil
}

func (client *OnlineConfClient) CreateEmptyNode(key string, skipAlreadyExist bool) error {
	params := map[string]string{
		"summary":      "",
		"description":  "",
		"notification": "",
		"mime":         "application/x-null",
		"data":         "",
		"comment":      "init key",
	}

	statusCode, result, err := client.request(key, http.MethodPost, params)
	log.Printf("init key %s, status: %+v, result: %+v, err: %+v\n", key, statusCode, result, err)
	if err != nil {
		return err
	}

	var response OnlineConfResponse
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		if response.Error == "AlreadyExists" && skipAlreadyExist {
			return nil
		}
		return fmt.Errorf("create empty node failure...status: %v, result: %v", statusCode, result)
	}
	return err
}

func (client *OnlineConfClient) CreateNode(item OnlineConfItem) error {

	params := map[string]string{
		"summary":      "",
		"description":  "",
		"notification": "",
		"mime":         item.Type,
		"data":         item.Value,
		"comment":      "init value",
	}

	log.Printf("creation key: %+v\n", item.Key)

	statusCode, result, err := client.request(item.Key, http.MethodPost, params)
	log.Printf("POST status: %+v, result: %+v, err: %+v\n", statusCode, result, err)
	if err != nil {

		return err
	}

	if statusCode != http.StatusOK {
		err := fmt.Errorf("create node failure...status: %v, result: %v", statusCode, result)
		log.Printf("ERROR: err: %+v\n", err)
		return nil
	}
	return err
}

func (client *OnlineConfClient) DeleteNode(key string) error {

	statusCode, result, err := client.request(key, http.MethodGet, nil)
	if statusCode != http.StatusOK {
		log.Printf("GET status: %+v, result: %+v, err: %+v\n", statusCode, result, err)
		return nil
	}

	//{"name":"qwerty","path":"/revise/test/qwerty","data":"1,2,3,ayewr","mime":"application/x-list","summary":"","description":"","version":3,"mtime":"2023-07-09 00:50:29","num_children":0,"access_modified":false,"rw":true,"notification":"none","notification_modified":false,"children":[]},
	var response OnlineConfResponse
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		return err
	}

	params := map[string]string{
		"version": strconv.Itoa(response.Version),
		"comment": "autoremove value",
	}

	log.Printf("delete key: %+v\n", key)

	statusCode, result, err = client.request(key, http.MethodDelete, params)
	log.Printf("DELETE status: %+v, result: %+v, err: %+v\n", statusCode, result, err)
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("delete node failure...status: %v, result: %v", statusCode, result)
	}
	return err
}

func (client *OnlineConfClient) request(
	requestUrl string,
	method string,
	params map[string]string,
) (int, string, error) {
	requestParams := url.Values{}
	for paramKey, paramValue := range params {
		requestParams.Add(paramKey, paramValue)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Prepare HTTP Client
	httpClient := &http.Client{Transport: transport}
	var reader io.Reader = nil
	if method != http.MethodGet {
		reader = strings.NewReader(requestParams.Encode())
	}

	url := fmt.Sprintf("%s/%s", client.host, requestUrl)

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return 0, "", fmt.Errorf("can't create the http request...%s", err.Error())
	}

	if method == http.MethodGet {
		req.URL.RawQuery = requestParams.Encode()
	}

	// Set headers
	for headerName, headerValue := range client.headers {
		req.Header.Set(headerName, headerValue)

	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Call the request
	log.Printf("request url: %s\n", req.URL.String())
	res, err := httpClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("can't call the http request...%s", err.Error())
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, "", fmt.Errorf("can't read http response...%s", err.Error())
	}

	return res.StatusCode, string(bodyBytes), nil
}
