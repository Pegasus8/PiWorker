package types

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

var c = map[PWType][]PWType{
	Any: {
		Text,
		Time,
		Int,
		Float,
		Bool,
		Path,
		JSON,
		URL,
		Date,
		Time,
	},

	Text: {
		Time,
		Int,
		Float,
		Bool,
		Path,
		JSON,
		URL,
		Date,
		Time,
	},

	Int: {
		Any,
		Text,
	},

	Float: {
		Any,
		Text,
		Int,
	},

	Bool: {
		Any,
		Text,
		Int,
		Float,
	},

	Path: {
		Any,
		Text,
	},

	JSON: {
		Any,
		Text,
	},

	URL: {
		Any,
		Text,
	},

	Date: {
		Any,
		Text,
	},

	Time: {
		Any,
		Text,
	},
}

// CompatWith checks if a type is compatible with other type (`targetType`).
func (t PWType) CompatWith(typeTarget PWType) bool {
	if _, exists := c[t]; !exists {
		log.Panic().Err(fmt.Errorf("unrecognized type '%v'", t)).Caller(1).Msg("")
	}

	for _, v := range c[t] {
		if v == typeTarget {
			return true
		}
	}

	return false
}

// CompatList returns a map with the types and their compatible pairs.
func CompatList() map[PWType][]PWType {
	return c
}
