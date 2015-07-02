package lol

import (
	"io/ioutil"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

type TempHomer struct {
	envOverrides map[string]string
	envSaves     map[string]string
	tmpDir       string
}

func (h *TempHomer) SetUp() error {
	var err error
	h.tmpDir, err = ioutil.TempDir("", "go-deb.ddesk_test")
	if err != nil {
	}
	h.OverrideEnv("HOME", h.tmpDir)

	h.envSaves = make(map[string]string)
	for key, value := range h.envOverrides {
		h.envSaves[key] = os.Getenv(key)
		os.Setenv(key, value)
	}
	return nil
}

func (h *TempHomer) TearDown() error {
	for key, value := range h.envSaves {
		os.Setenv(key, value)
	}
	return os.RemoveAll(h.tmpDir)
}

func (h *TempHomer) OverrideEnv(key, value string) {
	if h.envOverrides == nil {
		h.envOverrides = make(map[string]string)
	}
	h.envOverrides[key] = value
}

type XdgAPIKeyStorerSuite struct {
	h TempHomer
}

var _ = Suite(&XdgAPIKeyStorerSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *XdgAPIKeyStorerSuite) SetUpSuite(c *C) {
	s.h.OverrideEnv("XDG_CONFIG_HOME", "")
	err := s.h.SetUp()
	c.Assert(err, IsNil, Commentf("Initialization error: %s"))

}

func (s *XdgAPIKeyStorerSuite) TearDownSuite(c *C) {
	err := s.h.TearDown()
	c.Assert(err, IsNil, Commentf("Cleanup error: %s", err))
}

func (s *XdgAPIKeyStorerSuite) TestKeyValidity(c *C) {
	// These are old valid key got from riot games, invalid now :-P
	validData := []string{
		"b7b2daaf-ef84-41a0-8d79-de45ff42d311",
		"8975503a-9e88-4252-b11e-fc649f2dbe0d",
	}

	for _, keyStr := range validData {
		k := APIKey(keyStr)
		err := k.Check()
		c.Check(err, IsNil, Commentf("Got unexpected error: %s", err))
	}

}
