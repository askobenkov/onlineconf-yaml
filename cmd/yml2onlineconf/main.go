package main

import (
	"flag"
	"fmt"
	"reflect"
	"regexp"

	"log"

	client "onlineconf-yaml/onlineconf"
	"onlineconf-yaml/yml/parser"
)

/*
go run cmd/yml2onlineconf/main.go -onlineConfURL https://onlineconf.local -importConfigFilepath ./importConfig.yml -headersFilepath ./headers.txt -mainNodeName importConfig -showParsedConfig -importParsedConfig
*/
func main() {

	onlineConfURL := flag.String("onlineConfURL", "https://onlineconf.local", "OnlineConf URL name")
	configFilepath := flag.String("importConfigFilepath", "", "import config filepath")
	headersFilepath := flag.String("headersFilepath", "", "file with raw browser headers")
	mainNodeName := flag.String("mainNodeName", "", "OnlineConf main node name")
	showParsedConfig := flag.Bool("showParsedConfig", false, "Show parsed config")
	importParsedConfig := flag.Bool("importParsedConfig", false, "Import parsed config to OnlineConf")
	updateIfExists := flag.Bool("updateIfExists", false, "Update node value if already exists")
	deleteParsedConfig := flag.Bool("deleteParsedConfig", false, "Delete config in OnlineConf")
	skipAlreadyExist := flag.Bool("skipAlreadyExist", false, "Skip already exist error")
	skipCreateNode := flag.Bool("skipCreateNode", false, "Skip create node")
	basicAuthKey := flag.String("basicAuthKey", "", "Basic autorization key (docker only)")
	comment := flag.String("comment", "", "Comment message")

	flag.Parse()

	if *configFilepath == "" {
		log.Fatal(fmt.Errorf("import filepath config is empty"))
	}

	*onlineConfURL = regexp.MustCompile(`/+$`).ReplaceAllString(*onlineConfURL, "")
	client, err := client.NewOnlineConfClient(
		fmt.Sprintf("%s/%s/%s", *onlineConfURL, client.URLPrefix, *mainNodeName),
		*headersFilepath,
		*basicAuthKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	data, err := parser.GetYMLConfig(*configFilepath)
	if err != nil {
		log.Fatal(err)
	}

	obj := reflect.ValueOf(&data)
	src := parser.WalkByYML(obj, "", false)

	if *showParsedConfig {
		for k, v := range src {
			fmt.Printf("%-50s (%-30s) : %v\n", k, v.Type, v.Value)
		}
	}

	nodeKeys := parser.GetParentNodeKeys(src)

	if *importParsedConfig {
		if !*skipCreateNode {
			log.Printf("client.CreateEmptyNode..")

			for _, key := range nodeKeys {
				err := client.SetComment(*comment).CreateEmptyNode(key, *skipAlreadyExist)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		for _, v := range src {

			err = client.SetComment(*comment).CreateNode(v, *updateIfExists, *skipAlreadyExist)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if *deleteParsedConfig {
		nodeKeys = parser.GetNodeKeysForDelete(src)
		log.Printf("delete =============> %+v\n", nodeKeys)

		for _, key := range nodeKeys {
			err := client.SetComment(*comment).DeleteNode(key)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
