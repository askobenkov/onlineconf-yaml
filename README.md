# onlineconf-yaml

## yml2onlineconf - utility for import yaml config to OnlineConf by admin web interface

Options:
* onlineConfURL - onlineconf web interface URL
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
yml2onlineconf -onlineConfURL https://onlineconf.local -importConfigFilepath ./importConfig.yml -headersFilepath ./headers.txt -mainNodeName importConfig -showParsedConfig -importParsedConfig
```

## yml2cdb - utility for convert yml config to cdb database

Options:
* ymlConfigFilepath - input filepath to yml config
* cdbConfigFilepath - output filepath to cdb database
* [showParsedConfig] - show parsed config

Run:
```
yml2cdb -ymlConfigFilepath ./config.yml -cdbConfigFilepath ./config.cdb -showParsedConfig
```
