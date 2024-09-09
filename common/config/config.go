package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Production                  bool             `json:"production"`
	AppName                     string           `json:"appName"`
	KubernetesServiceNameSuffix string           `json:"kubernetesServiceNameSuffix"`
	ListenPort                  string           `json:"listenPort"`
	Folders                     []string         `json:"folders"`
	Jaeger                      JaegerConfig     `json:"jaeger"`
	Consul                      ConsulConfig     `json:"consul"`
	GrpcServer                  GrpcServerConfig `json:"grpcServer"`
}

type JaegerConfig struct {
	JaegerEndpoint string `json:"jaegerEndpoint"`
	ServiceName    string `json:"serviceName"`
	ServiceVersion string `json:"serviceVersion"`
}

type GrpcServerConfig struct {
	Port              string `json:"port"`
	MaxConnectionIdle int    `json:"maxConnectionIdle"`
	MaxConnectionAge  int    `json:"maxConnectionAge"`
	Timeout           int    `json:"timeout"`
}

type ConsulConfig struct {
	Host string `json:"host"`
}

func MustLoadConfig(production bool, path string) *Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("config.dev")
	if production {
		viper.SetConfigName("config.prod")
	}
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	config := &Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %s", err))
	}

	config.Production = production

	return config
}
