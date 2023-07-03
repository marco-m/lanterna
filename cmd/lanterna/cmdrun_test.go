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

var cfg = Args{
	log:        zerolog.Logger{},
	ConfigPath: "testdata/config.json",
}

func TestCmdRunHappyPathMock(t *testing.T) {
	collect := collectMock{ips: []string{"1.2.3.4"}}
	post := postMock{}

	err := cmdRun(cfg, collect.collect, post.postJSON)

	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.StringContains(post.url, "http://mango.example&threadKey="))
	qt.Assert(t, qt.StringContains(post.msg["text"], "IP addresses:\n    1.2.3.4"))
}

func TestCmdRunNoAddressFoundSentAsWarningMock(t *testing.T) {
	collect := collectMock{}
	pMock := postMock{}

	err := cmdRun(cfg, collect.collect, pMock.postJSON)

	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.StringContains(pMock.msg["text"], "IP addresses:\n    WARNING: none found"))
}
