package errors

import (
	"strings"
	"sync"
)

var (
	stringBuilderPool = sync.Pool{
		New: func() any {
			return &strings.Builder{}
		},
	}
)
