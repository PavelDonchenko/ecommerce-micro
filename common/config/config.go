package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Production                  bool                  `json:"production"`
	AppName                     string                `json:"appName"`
	KubernetesServiceNameSuffix string                `json:"kubernetesServiceNameSuffix"`
	ListenPort                  string                `json:"listenPort"`
	Folders                     []string              `json:"folders"`
	Jaeger                      JaegerConfig          `json:"jaeger"`
	Consul                      ConsulConfig          `json:"consul"`
	GrpcServer                  GrpcServerConfig      `json:"grpcServer"`
	Certificates                CertificatesConfig    `json:"certificates"`
	SecurityKeys                SecurityKeysConfig    `json:"securityKeys"`
	SecurityRSAKeys             SecurityRSAKeysConfig `json:"securityRSAKeys"`
	EmailService                EmailServiceConfig    `json:"emailService"`
	SecondsToReloadServicesName int                   `json:"secondsToReloadServicesName"`
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

type EmailServiceConfig struct {
	ServiceName string `json:"serviceName"`
	Host        string `json:"host"`
}

type CertificatesConfig struct {
	FolderName                    string `json:"foldername"`
	FileNameCert                  string `json:"filenamecert"`
	FileNameKey                   string `json:"filenamekey"`
	HashPermissionEndPoint        string `json:"hashPermissionEndPoint"`
	PasswordPermissionEndPoint    string `json:"passwordPermissionEndPoint"`
	ServiceName                   string `json:"serviceName"`
	APIPathCertificateCA          string `json:"apiPathCertificateCA"`
	EndPointGetCertificateCA      string `json:"endPointGetCertificateCA"`
	APIPathCertificateHost        string `json:"apiPathCertificateHost"`
	EndPointGetCertificateHost    string `json:"endPointGetCertificateHost"`
	APIPathCertificateHostKey     string `json:"apiPathCertificateHostKey"`
	EndPointGetCertificateHostKey string `json:"endPointGetCertificateHostKey"`
	MinutesToReloadCertificate    int    `json:"minutesToReloadCertificate"`
}

type SecurityKeysConfig struct {
	DaysToExpireKeys            int    `json:"daysToExpireKeys"`
	MinutesToRefreshPrivateKeys int    `json:"minutesToRefreshPrivateKeys"`
	MinutesToRefreshPublicKeys  int    `json:"minutesToRefreshPublicKeys"`
	SavePublicKeyToFile         bool   `json:"savePublicKeyToFile"`
	FileECPPublicKey            string `json:"fileECPPublicKey"`
	ServiceName                 string `json:"serviceName"`
	APIPathPublicKeys           string `json:"apiPathPublicKeys"`
	EndPointGetPublicKeys       string `json:"endPointGetPublicKeys"`
}

type SecurityRSAKeysConfig struct {
	DaysToExpireRSAKeys            int    `json:"daysToExpireRSAKeys"`
	MinutesToRefreshRSAPrivateKeys int    `json:"minutesToRefreshRSAPrivateKeys"`
	MinutesToRefreshRSAPublicKeys  int    `json:"minutesToRefreshRSAPublicKeys"`
	ServiceName                    string `json:"serviceName"`
	APIPathRSAPublicKeys           string `json:"apiPathRSAPublicKeys"`
	EndPointGetRSAPublicKeys       string `json:"endPointGetRSAPublicKeys"`
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
