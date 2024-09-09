package consul

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PavelDonchenko/template/MICROSERVICES/ecommerce-micro/common/config"
	consul "github.com/hashicorp/consul/api"
)

type ConsulParse string

const (
	CertificatesAndSecurityKeys ConsulParse = "CertificatesAndSecurityKeys"
	SecurityRSAKeys             ConsulParse = "SecurityRSAKeys"
	EmailService                ConsulParse = "EmailService"
)

func NewConsulClient(cfg *config.Config) (*consul.Client, string, error) {
	consulConfig := consul.DefaultConfig()
	consulConfig.Address = cfg.Consul.Host

	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, "", err
	}

	serviceID, err := register(cfg, consulClient)
	if err != nil {
		return nil, "", err
	}

	return consulClient, serviceID, nil
}

func register(cfg *config.Config, client *consul.Client) (string, error) {
	address := hostname()

	k8s, _ := strconv.ParseBool(os.Getenv("kubernetes"))
	if k8s {
		address = fmt.Sprintf("%s-%s", cfg.AppName, cfg.KubernetesServiceNameSuffix)
	}

	port, err := strconv.Atoi(strings.Split(cfg.ListenPort, ":")[1])
	if port == 0 || err != nil {
		return "", err
	}

	if len(strings.TrimSpace(cfg.GrpcServer.Port)) > 0 {
		port, err = strconv.Atoi(strings.Split(cfg.GrpcServer.Port, ":")[1])
		if err != nil {
			return "", err
		}
	}

	serviceID := fmt.Sprintf("%s-%s:%v", cfg.AppName, address, port)

	httpCheck := fmt.Sprintf("https://%s:%v/healthy", address, port)
	fmt.Println("Consul httpCheck:", httpCheck)

	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    cfg.AppName,
		Port:    port,
		Address: address,
		Check: &consul.AgentServiceCheck{
			CheckID:                        serviceID,
			Name:                           fmt.Sprintf("Service %s check", cfg.AppName),
			HTTP:                           httpCheck,
			TLSSkipVerify:                  true,
			Interval:                       "10s",
			Timeout:                        "30s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Printf("failed consul to register service: %s:%v ", address, port)
		return "", err
	}

	log.Printf("successfully consul register service: %s:%v", address, port)

	return serviceID, nil
}

func hostname() string {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	return hostName
}
