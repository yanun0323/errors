package errors

import (
	"fmt"
	"testing"
)

func TestReplaceFormatError(t *testing.T) {
	testCases := []struct {
		desc     string
		format   string
		arg      []any
		expected string
	}{
		{
			"Error with single %+v",
			"there is an %s error, err: %+v",
			[]any{"hello", New("root error")},
			"there is an %s error, err: %s",
		},
		{
			"Error with single %v",
			"there is an %s error, err: %v",
			[]any{"hello", New("root error")},
			"there is an %s error, err: %s",
		},
		{
			"Error with single %#v",
			"there is an %s error, err: %#v",
			[]any{"hello", New("root error")},
			"there is an %s error, err: %s",
		},
		{
			"Error with multiple %+v",
			"there is an %+v error, err: %+v",
			[]any{"hello", New("root error")},
			"there is an %+v error, err: %s",
		},
		{
			"Error with multiple %v",
			"there is an %v error, err: %v",
			[]any{"hello", New("root error")},
			"there is an %v error, err: %s",
		},
		{
			"Error with multiple %#v",
			"there is an %#v error, err: %#v",
			[]any{"hello", New("root error")},
			"there is an %#v error, err: %s",
		},
		{
			"fmt.Error with single %+v",
			"there is an %s error, err: %+v",
			[]any{"hello", fmt.Errorf("root error")},
			"there is an %s error, err: %+v",
		},
		{
			"fmt.Error with multiple %+v",
			"there is an %+v error, err: %+v",
			[]any{"hello", fmt.Errorf("root error")},
			"there is an %+v error, err: %+v",
		},
		{
			"complex error",
			"this is %%%% error1: %v, this is error2: %+v, %%%% this is float: %f%% this is error3: %#v, this is error4: %v",
			[]any{
				New("error1"),
				fmt.Errorf("error2"),
				0.323,
				Wrap(New("error3 root"), "error3"),
				fmt.Errorf("error4, err: %+v", "no!!"),
			},
			"this is %%%% error1: %s, this is error2: %+v, %%%% this is float: %f%% this is error3: %s, this is error4: %v",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := replaceFormatError(tc.format, tc.arg...)
			if tc.expected != got {
				t.Errorf("expected: '%s', but got '%s'", tc.expected, got)
			}
		})
	}
}
