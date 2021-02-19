package common

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Pastel	PastelConfiguration		`yaml:"pastel"`
	Storage StorageConfiguration	`yaml:"storage"`
}

type PastelConfiguration struct {
	DataDir string `yaml:"data-dir"`
}

type StorageConfiguration struct {
	BootstrapNodes []string `yaml:"bootstrapNodes"`
	RpcHost string 			`yaml:"rpc-host"`
	RpcPort int 			`yaml:"rpc-port"`
	Port int 				`yaml:"p2p-port"`
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
