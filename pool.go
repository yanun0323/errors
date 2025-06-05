package errors

import (
	"strings"
	"sync"
)

var (
	dataPool = sync.Pool{
		New: func() any {
			return make(map[string]any, 4)
		},
	}

	stringBuilderPool = sync.Pool{
		New: func() any {
			return &strings.Builder{}
		},
	}
)
