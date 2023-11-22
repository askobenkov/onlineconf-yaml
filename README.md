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
* [basicAuthKey] - Basic autorization key (docker only)

Run:
```
go run onlineconf.go -onlineConfUrl https://onlineconf.local -exportConfigFilepath ./exportConfig.yml -headersFilepath ./headers.txt -mainNodeName exportConfig -showParsedConfig -exportParsedConfig
```
