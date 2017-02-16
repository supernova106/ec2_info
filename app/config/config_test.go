package config

import (
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ConfigSuite struct{}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) TestLoad(c *C) {
	os.Setenv("OAUTH2_URL", "https://oauth2-staging.razerapi.com")
	os.Setenv("GIN_ENV", "testing")

	defer os.Setenv("OAUTH2_URL", "")
	defer os.Setenv("GIN_ENV", "")

	config, err := Load(".env")
	c.Assert(err, IsNil)
	c.Assert(config.MysqlUrl, Equals, "https://oauth2-staging.razerapi.com")
	c.Assert(config.GinEnv, Equals, "testing")
}

