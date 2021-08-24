// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/util/flagext/stringslicecsv_test.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package flagext

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func Test_StringSliceCSV(t *testing.T) {
	type TestStruct struct {
		CSV StringSliceCSV `yaml:"csv"`
	}

	var testStruct TestStruct
	s := "a,b,c,d"
	assert.Nil(t, testStruct.CSV.Set(s))

	assert.Equal(t, []string{"a", "b", "c", "d"}, []string(testStruct.CSV))
	assert.Equal(t, s, testStruct.CSV.String())

	expected := []byte(`csv: a,b,c,d
`)

	actual, err := yaml.Marshal(testStruct)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	var testStruct2 TestStruct

	err = yaml.Unmarshal(expected, &testStruct2)
	assert.Nil(t, err)
	assert.Equal(t, testStruct, testStruct2)
}