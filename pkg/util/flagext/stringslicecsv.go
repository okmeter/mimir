// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/util/flagext/stringslicecsv.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package flagext

import "strings"

// StringSliceCSV is a slice of strings that is parsed from a comma-separated string
// It implements flag.Value and yaml Marshalers
type StringSliceCSV []string

// String implements flag.Value
func (v StringSliceCSV) String() string {
	return strings.Join(v, ",")
}

// Set implements flag.Value
func (v *StringSliceCSV) Set(s string) error {
	*v = strings.Split(s, ",")
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (v *StringSliceCSV) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	return v.Set(s)
}

// MarshalYAML implements yaml.Marshaler.
func (v StringSliceCSV) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}