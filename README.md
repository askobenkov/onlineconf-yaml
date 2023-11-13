# onlineconf-yaml
Utility for import yaml config to onlineconf by admin web interface

Options:
* onlineConfUrl - onlineconf web interface URL
* exportConfigFilepath - filepath to yaml config
* headersFilepath - filepath to http headers
* mainNodeName - name of the node where the config will be imported
* [showParsedConfig] - show parsed config
* [exportParsedConfig] - export parsed config
* [deleteParsedConfig] - delete parsed config
* [skipAlreadyExist] - skip already exists error
* [skipCreateNode] - skip node creating

Run:
```
go run onlineconf.go -onlineConfUrl https://onlineconf.dev.dmr/config -exportConfigFilepath ./Revise1.yml -headersFilepath ./headers.txt -mainNodeName revise -showParsedConfig
```
