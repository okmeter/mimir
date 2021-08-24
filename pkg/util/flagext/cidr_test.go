// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/util/flagext/cidr_test.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package flagext

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type TestStruct struct {
	CIDRs CIDRSliceCSV `yaml:"cidrs" json:"cidrs"`
}

func Test_CIDRSliceCSV_YamlMarshalling(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"should marshal empty config": {
			input:    "cidrs: \"\"\n",
			expected: nil,
		},
		"should marshal single value": {
			input:    "cidrs: 127.0.0.1/32\n",
			expected: []string{"127.0.0.1/32"},
		},
		"should marshal multiple comma-separated values": {
			input:    "cidrs: 127.0.0.1/32,10.0.10.0/28,fdf8:f53b:82e4::/100,192.168.0.0/20\n",
			expected: []string{"127.0.0.1/32", "10.0.10.0/28", "fdf8:f53b:82e4::/100", "192.168.0.0/20"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Unmarshal.
			actual := TestStruct{}
			err := yaml.Unmarshal([]byte(tc.input), &actual)
			assert.NoError(t, err)

			assert.Len(t, actual.CIDRs, len(tc.expected))
			for idx, cidr := range actual.CIDRs {
				assert.Equal(t, tc.expected[idx], cidr.String())
			}

			// Marshal.
			out, err := yaml.Marshal(actual)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, string(out))
		})
	}
}

func Test_CIDRSliceCSV_JSONMarshalling(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"should marshal empty config": {
			input:    `{"cidrs":""}`,
			expected: nil,
		},
		"should marshal single value": {
			input:    `{"cidrs":"127.0.0.1/32"}`,
			expected: []string{"127.0.0.1/32"},
		},
		"should marshal multiple comma-separated values": {
			input:    `{"cidrs":"127.0.0.1/32,10.0.10.0/28,fdf8:f53b:82e4::/100,192.168.0.0/20"}`,
			expected: []string{"127.0.0.1/32", "10.0.10.0/28", "fdf8:f53b:82e4::/100", "192.168.0.0/20"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Unmarshal.
			actual := TestStruct{}
			err := json.Unmarshal([]byte(tc.input), &actual)
			assert.NoError(t, err)

			assert.Len(t, actual.CIDRs, len(tc.expected))
			for idx, cidr := range actual.CIDRs {
				assert.Equal(t, tc.expected[idx], cidr.String())
			}

			// Marshal.
			out, err := json.Marshal(actual)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, string(out))
		})
	}
}