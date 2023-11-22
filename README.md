# onlineconf-yaml

Utility for import yaml config to OnlineConf by admin web interface

Options:
* onlineConfUrl - onlineconf web interface URL
* importConfigFilepath - filepath to yaml config
* headersFilepath - filepath to http headers
* mainNodeName - name of the node where the config will be imported
* [showParsedConfig] - show parsed config
* [importParsedConfig] - import parsed config
* [deleteParsedConfig] - delete parsed config
* [skipAlreadyExist] - skip already exists error
* [skipCreateNode] - skip node creating
* [basicAuthKey] - Basic autorization key (docker only)

Run:
```
go run onlineconf.go -onlineConfUrl https://onlineconf.local -importConfigFilepath ./importConfig.yml -headersFilepath ./headers.txt -mainNodeName importConfig -showParsedConfig -importParsedConfig
```
