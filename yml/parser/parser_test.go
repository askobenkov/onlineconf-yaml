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
    R:
      "1001-1001-1001":
        KEY1:
          - fee: 4
            pmtype: FEEKEY1
            subject: "Simple fee"
  volatile:
    VR:
      "1001-1001-1001":
        KEY1:
          - fee: 4
            pmtype: FEEKEY1
            subject: "Simple fee"
      "1006-1001-1001":
        KEY3:
          - fee: 2.15
            pmtype: FEEKEY3
            subject: "Simple fee"
        KEY2:
          - fee: 2.15
            pmtype: FEEKEY2
            subject: "Simple fee"
  parent_nested:
    sub_parent01: {}
    reg-exp-.*?@part01\.part02\.part03\.ru: {}
`)

	require.NoError(t, err)

	data, err := GetYMLConfig(cfgFilepath)
	require.NoError(t, err)

	expected := map[string]OnlineConfItem{
		"fee.": {
			Key:   "fee.",
			Value: "[\"common\",\"parent_nested\",\"volatile\"]",
			Type:  "application/x-yaml",
		},
		"fee/common.": {
			Key:   "fee/common.",
			Value: "[\"R\"]",
			Type:  "application/x-yaml",
		},
		"fee/common/R.": {
			Key:   "fee/common/R.",
			Value: "[\"1001-1001-1001\"]",
			Type:  "application/x-yaml",
		},
		"fee/common/R/1001-1001-1001.": {
			Key:   "fee/common/R/1001-1001-1001.",
			Value: "[\"KEY1\"]",
			Type:  "application/x-yaml",
		},
		"fee/common/R/1001-1001-1001/KEY1": {
			Key:   "fee/common/R/1001-1001-1001/KEY1",
			Value: "- fee: 4\n  pmtype: FEEKEY1\n  subject: Simple fee\n",
			Type:  "application/x-yaml",
		},
		"fee/parent_nested.": {
			Key:   "fee/parent_nested.",
			Value: "[\"reg-exp-.*?@part01\\\\.part02\\\\.part03\\\\.ru\",\"sub_parent01\"]",
			Type:  "application/x-yaml",
		},
		"fee/parent_nested/sub_parent01": {
			Key:   "fee/parent_nested/sub_parent01",
			Value: "{}",
			Type:  "application/x-yaml",
		},
		"fee/parent_nested/reg-exp-.*?@part01\\.part02\\.part03\\.ru": {
			Key:   "fee/parent_nested/reg-exp-.*?@part01\\.part02\\.part03\\.ru",
			Value: "{}",
			Type:  "application/x-yaml",
		},
		"fee/volatile.": {
			Key:   "fee/volatile.",
			Value: "[\"VR\"]",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VR.": {
			Key:   "fee/volatile/VR.",
			Value: "[\"1001-1001-1001\",\"1006-1001-1001\"]",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VR/1001-1001-1001.": {
			Key:   "fee/volatile/VR/1001-1001-1001.",
			Value: "[\"KEY1\"]",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VR/1001-1001-1001/KEY1": {
			Key:   "fee/volatile/VR/1001-1001-1001/KEY1",
			Value: "- fee: 4\n  pmtype: FEEKEY1\n  subject: Simple fee\n",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VR/1006-1001-1001.": {
			Key:   "fee/volatile/VR/1006-1001-1001.",
			Value: "[\"KEY2\",\"KEY3\"]",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VR/1006-1001-1001/KEY3": {
			Key:   "fee/volatile/VR/1006-1001-1001/KEY3",
			Value: "- fee: 2.15\n  pmtype: FEEKEY3\n  subject: Simple fee\n",
			Type:  "application/x-yaml",
		},
		"fee/volatile/VR/1006-1001-1001/KEY2": {
			Key:   "fee/volatile/VR/1006-1001-1001/KEY2",
			Value: "- fee: 2.15\n  pmtype: FEEKEY2\n  subject: Simple fee\n",
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
