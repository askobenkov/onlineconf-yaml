package client

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"onlineconf-yaml/yml/parser"
	"os"
	"strconv"
	"strings"
)

// URLPrefix prefix url constant
const URLPrefix = "config"

// OnlineConfClient onlineconf client
type OnlineConfClient struct {
	host    string
	headers map[string]string
	comment string
}

// OnlineConfResponse onlineconf response
type OnlineConfResponse struct {
	Error   string
	Version int
	Message string
}

// NewOnlineConfClient create onlineconf client
func NewOnlineConfClient(
	host string,
	filepathHeader string,
	basicAuthKey string,
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
	} else {
		client.headers = map[string]string{
			"X-Requested-With": "XMLHttpRequest",
			"Authorization":    fmt.Sprintf("Basic %s", basicAuthKey),
		}
	}

	return client, nil
}

// SetComment set comment
func (client *OnlineConfClient) SetComment(msg string) *OnlineConfClient {
	client.comment = msg
	return client
}

// GetHeaders return headers
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
		if len(list) != 2 {
			continue
		}
		headers[strings.TrimSpace(list[0])] = strings.TrimSpace(list[1])
	}
	if err := scanner.Err(); err != nil {
		return headers, err
	}
	return headers, nil
}

// CreateEmptyNode creating empty node
func (client *OnlineConfClient) CreateEmptyNode(key string, skipAlreadyExist bool) error {
	params := map[string]string{
		"summary":      "",
		"description":  "",
		"notification": "",
		"mime":         "application/x-null",
		"data":         "",
		"comment":      "init key",
	}
	if client.comment != "" {
		params["comment"] = client.comment
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

// CreateNode create node
func (client *OnlineConfClient) CreateNode(item parser.OnlineConfItem, updateIfExists bool, skipAlreadyExist bool) error {

	params := map[string]string{
		"summary":      "",
		"description":  "",
		"notification": "",
		"mime":         item.Type,
		"data":         item.Value,
		"comment":      "init value",
	}
	if client.comment != "" {
		params["comment"] = client.comment
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

		if statusCode != http.StatusBadRequest {
			return err
		}

		if updateIfExists {

			statusCode, result, err = client.request(item.Key, http.MethodGet, nil)
			if statusCode != http.StatusOK {
				log.Printf("GET status: %+v, result: %+v, err: %+v\n", statusCode, result, err)
				return nil
			}

			var response OnlineConfResponse
			err = json.Unmarshal([]byte(result), &response)
			if err != nil {
				return err
			}
			params["version"] = strconv.Itoa(response.Version)

			log.Printf("update key: %+v\n", item.Key)

			statusCode, result, err := client.request(item.Key, http.MethodPost, params)
			log.Printf("POST status: %+v, result: %+v, err: %+v\n", statusCode, result, err)
			if err != nil {
				return err
			}
			return nil
		}

		if statusCode == http.StatusBadRequest && skipAlreadyExist {
			return nil
		}

		return err
	}
	return err
}

// DeleteNode delete node
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

	if client.comment != "" {
		params["comment"] = client.comment
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
	requestURL string,
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

	url := fmt.Sprintf("%s/%s", client.host, requestURL)

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

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, "", fmt.Errorf("can't read http response...%s", err.Error())
	}

	return res.StatusCode, string(bodyBytes), nil
}
