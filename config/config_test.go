package config

import (
	. "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type Tests struct{}

var _ = Suite(&Tests{})

func (*Tests) TestLoadConfig(c *C) {
	configRoot = "test"
	def := filepath.FromSlash(configRoot + "/defaults.json")
	path := filepath.FromSlash(configRoot + "/config.json")
	standard, err := ioutil.ReadFile(def)
	c.Assert(err, IsNil)

	// Config file does not exist
	c.Assert(LoadConfig(), ErrorMatches, "Error reading configuration file.*")

	c.Assert(ioutil.WriteFile(path, standard, 0600), IsNil)
	defer func() {
		c.Assert(os.Remove(path), IsNil)
	}()

	c.Assert(LoadConfig(), IsNil)
	stdConfig := &ServerConfigs{}
	stdConfig.Posts.Salt = "LALALALALALALALALALALALALALALALALALALALA"
	c.Assert(config, DeepEquals, stdConfig)
	c.Assert(hash, Equals, "eeba38176564a577")
}

func (*Tests) TestSettingAndGetting(c *C) {
	conf := ServerConfigs{}
	conf.Boards.Enabled = []string{"a", "l", "k"}
	Set(conf)
	c.Assert(Get(), DeepEquals, &conf)
}

func (*Tests) TestSettingAndGettingClient(c *C) {
	std := []byte{1, 2, 3}
	hash := "foo"
	SetClient(std, hash)
	json, jsonHash := GetClient()
	c.Assert(json, DeepEquals, std)
	c.Assert(jsonHash, Equals, hash)
}
