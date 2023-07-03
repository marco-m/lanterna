package main

import (
	"testing"

	"github.com/go-quicktest/qt"
	"github.com/rs/zerolog"
)

type collectMock struct {
	ips []string
}

func (cm *collectMock) collect(log zerolog.Logger) ([]string, error) {
	return cm.ips, nil
}

type postMock struct {
	url string
	msg map[string]string
}

func (pm *postMock) postJSON(url string, msg map[string]string) error {
	pm.url = url
	pm.msg = msg
	return nil
}

var config = configuration{
	Sinks: []sink{
		{
			Name: "banana",
			Type: "gchat",
			URL:  "http://mango.example",
		},
	},
}

func TestRunHandleHappyPathMock(t *testing.T) {
	collect := collectMock{ips: []string{"1.2.3.4"}}
	post := postMock{}

	err := runHandle(Args{}, config, "margherita", collect.collect, post.postJSON)

	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.StringContains(post.url, "http://mango.example&threadKey="))
	qt.Assert(t, qt.StringContains(post.msg["text"], "IP addresses:\n    1.2.3.4"))
}

func TestRunHandleNoAddressFoundSentAsWarningMock(t *testing.T) {
	collect := collectMock{}
	pMock := postMock{}

	err := runHandle(Args{}, config, "margherita", collect.collect, pMock.postJSON)

	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.StringContains(pMock.msg["text"], "IP addresses:\n    WARNING: none found"))
}
