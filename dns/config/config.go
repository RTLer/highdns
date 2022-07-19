package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Upstream struct {
	Suffix string   `yaml:"suffix"`
	Hosts  []string `yaml:"hosts"`
}
type Server struct {
	Address string `yaml:"address"`
	Net     string `yaml:"net"`
}
type Conf struct {
	Server   Server     `yaml:"server"`
	Upstream []Upstream `yaml:"upstream"`
	Main     []string   `yaml:"main"`
}

var Config = &Conf{}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.Unmarshal(Config)
}

func (c *Conf) GetHosts(domain string) Upstream {
	for _, u := range c.Upstream {
		if strings.HasSuffix(domain[:len(domain)-1], u.Suffix) {
			return u
		}
	}
	return Upstream{
		Hosts: c.Main,
	}
}
