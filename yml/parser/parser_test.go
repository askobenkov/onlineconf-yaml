package parser

import (
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {

	cfgFilepath, err := writeYMLConfig(`
fee:
  common:
    RUR:
      "1001-9001-3049-1111":
        TOBLRBNK:
          - fee: 4
            pmtype: FEETOBLRBNK
            subject: "Simple fee"
  volatile:
    VRUR:
      "1001-9001-3049-1111":
        TOBLRBNK:
          - fee: 4
            pmtype: FEETOBLRBNK
            subject: "Simple fee"
      "1006-0001-5761-2222":
        TOSBRTR1:
          - fee: 2.15
            pmtype: FEETOSBRTR1
            subject: "Simple fee"
        TOSBRTR2:
          - fee: 2.15
            pmtype: FEETOSBRTR2
            subject: "Simple fee"
`)

	require.NoError(t, err)

	data, err := GetYMLConfig(cfgFilepath)
	require.NoError(t, err)

	expected := map[string]OnlineConfItem{
		"fee.": {
			Key:   "fee.",
			Value: "[\"common\",\"volatile\"]",
			Type:  "application/x-yaml",
		},
		"fee/common.": {
			Key:   "fee/common.",
			Value: "[\"RUR\"]",
			Type:  "application/x-yaml",
		},
		"fee/common/RUR.": {
			Key:   "fee/common/RUR.",
			Value: "[\"1001-9001-3049-1111\"]",
			Type:  "application/x-yaml",
		},
		"fee/common/RUR/1001-9001-3049-1111.": {
			Key:   "fee/common/RUR/1001-9001-3049-1111.",
			Value: "[\"TOBLRBNK\"]",
			Type:  "application/x-yaml",
		},
		"fee/common/RUR/1001-9001-3049-1111/TOBLRBNK": {
			Key:   "fee/common/RUR/1001-9001-3049-1111/TOBLRBNK",
			Value: "- fee: 4\n  pmtype: FEETOBLRBNK\n  subject: Simple fee\n",
			Type:  "application/x-yaml",
		},
		"fee/volatile.": {
			Key:   "fee/volatile.",
			Value: "[\"VRUR\"]",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VRUR.": {
			Key:   "fee/volatile/VRUR.",
			Value: "[\"1001-9001-3049-1111\",\"1006-0001-5761-2222\"]",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VRUR/1001-9001-3049-1111.": {
			Key:   "fee/volatile/VRUR/1001-9001-3049-1111.",
			Value: "[\"TOBLRBNK\"]",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VRUR/1001-9001-3049-1111/TOBLRBNK": {
			Key:   "fee/volatile/VRUR/1001-9001-3049-1111/TOBLRBNK",
			Value: "- fee: 4\n  pmtype: FEETOBLRBNK\n  subject: Simple fee\n",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VRUR/1006-0001-5761-2222.": {
			Key:   "fee/volatile/VRUR/1006-0001-5761-2222.",
			Value: "[\"TOSBRTR1\",\"TOSBRTR2\"]",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VRUR/1006-0001-5761-2222/TOSBRTR1": {
			Key:   "fee/volatile/VRUR/1006-0001-5761-2222/TOSBRTR1",
			Value: "- fee: 2.15\n  pmtype: FEETOSBRTR1\n  subject: Simple fee\n",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VRUR/1006-0001-5761-2222/TOSBRTR2": {
			Key:   "fee/volatile/VRUR/1006-0001-5761-2222/TOSBRTR2",
			Value: "- fee: 2.15\n  pmtype: FEETOSBRTR2\n  subject: Simple fee\n",
			Type:  "application/x-yaml",
		},
	}

	obj := reflect.ValueOf(&data)
	src := WalkByYML(obj, "", true)

	assert.Equal(t, true, reflect.DeepEqual(src, expected), "struct equals")
}

func writeYMLConfig(content string) (string, error) {
	f, err := ioutil.TempFile("", "testOnlineConf")
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(f, content)
	if err != nil {
		return "", err
	}
	return f.Name(), err
}
