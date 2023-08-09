package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-quicktest/qt"
)

var testCases = []struct {
	name string
	url  string
	want string
}{
	{
		"with params",
		"https://example.com/x?key=A&token=B",
		"https://example.com/x?REDACTED",
	},
	{
		"no params",
		"https://example.com/x",
		"https://example.com/x",
	},
	{
		"with user/pw and params",
		"https://user:pass@example.com/x?key=A&token=B",
		"https://REDACTED:REDACTED@example.com/x?REDACTED",
	},
}

func TestRedactURL(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			urlo, err := url.Parse(tc.url)
			qt.Assert(t, qt.IsNil(err))

			redacted := RedactURL(urlo)

			qt.Assert(t, qt.Equals(redacted.String(), tc.want))
		})
	}
}

func TestRedactURLString(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redacted := RedactURLString(tc.url)

			qt.Assert(t, qt.Equals(redacted, tc.want))
		})
	}
}

func TestRedactErrorURL(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			urlErr := url.Error{Op: http.MethodGet, URL: tc.url, Err: errors.New("banana")}

			redacted := RedactErrorURL(&urlErr)

			want := fmt.Sprintf("%s %q: banana", http.MethodGet, tc.want)
			qt.Assert(t, qt.Equals(redacted.Error(), want))
		})
	}
	// in case the error is not url.Error:
	fooErr := errors.New("foo")
	redacted := RedactErrorURL(fooErr)
	qt.Assert(t, qt.ErrorIs(redacted, fooErr))
}
