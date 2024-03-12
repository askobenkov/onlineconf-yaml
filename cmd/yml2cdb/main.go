package main

import (
	"flag"
	"fmt"
	"reflect"

	"log"

	"onlineconf-yaml/cdb"
	"onlineconf-yaml/yml/parser"
)

/*
go run cmd/yml2cdb/main.go -ymlConfigFilepath ./importConfig.yml -cdbConfigFilepath ./importConfig.cdb -showParsedConfig
*/

func main() {

	ymlConfigFilepath := flag.String("ymlConfigFilepath", "", "yml input config filepath")
	cdbConfigFilepath := flag.String("cdbConfigFilepath", "", "cdb output config filepath")
	showParsedConfig := flag.Bool("showParsedConfig", false, "Show parsed config")

	flag.Parse()

	if *ymlConfigFilepath == "" {
		log.Fatal(fmt.Errorf("input filepath config is empty"))
	}

	if *cdbConfigFilepath == "" {
		log.Fatal(fmt.Errorf("output filepath config is empty"))
	}

	data, err := parser.GetYMLConfig(*ymlConfigFilepath)
	if err != nil {
		log.Fatal(err)
	}

	obj := reflect.ValueOf(&data)
	src := parser.WalkByYML(obj, "", true)

	params := make([]cdb.WriteItem, len(src))
	for k, v := range src {
		if *showParsedConfig {
			log.Printf("%-50s (%-30s) : %v\n", k, v.Type, v.Value)
		}

		params = append(params, cdb.WriteItem{
			Path:  k,
			Value: v.Value,
			Tp:    v.Type,
		})
	}

	err = cdb.Write(*cdbConfigFilepath, params)
	if err != nil {
		log.Fatal(err)
	}
}
