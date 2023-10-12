package save_test

import (
	"errors"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com1",
		},
		{
			name:      "Empty url",
			alias:     "some_alias",
			url:       "",
			respError: "failed url is a required field",
		},
		{
			name:      "Invalid url",
			alias:     "some_alias1",
			url:       "some invalid url",
			respError: "failed url is a valid url",
		},
		{
			name:      "Save url error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
		},
	}

	for _, tc := range cases {
		tc := tc
		t

	}
}
