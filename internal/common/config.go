package common

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Pastel PastelConfiguration     `yaml:"pastel"`
	REST   RESTServerConfiguration `yaml:"rest"`
	P2P    P2PConfiguration        `yaml:"p2p"`
}

type PastelConfiguration struct {
	DataDir string `yaml:"data-dir"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
	Pwd     string `yaml:"pwd"`
}

type RESTServerConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type P2PConfiguration struct {
	Port int `yaml:"port"`
}

func (c *Config) LoadConfig(configFile string) error {

	viper.SetConfigName(configFile) // config file name without extension
	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %s", err)
	}

	err := viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	return nil
}
