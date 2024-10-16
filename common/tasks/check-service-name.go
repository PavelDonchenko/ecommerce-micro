package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	consul "github.com/hashicorp/consul/api"

	"github.com/PavelDonchenko/ecommerce-micro/common/config"
	common_consul "github.com/PavelDonchenko/ecommerce-micro/common/consul"
	trace "github.com/PavelDonchenko/ecommerce-micro/common/trace/otel"
)

func ReloadServiceName(ctx context.Context,
	config *config.Config,
	consulClient *consul.Client,
	serviceName string,
	consulParse common_consul.Parse,
	servicesNameDone chan bool) {
	ticker := time.NewTicker(2500 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_, span := trace.NewSpan(ctx, "checkServiceNameTask.ReloadServiceName")
				defer span.End()

				services, _, err := consulClient.Catalog().Service(serviceName, "", nil)
				if err != nil {
					fmt.Printf("failed to refresh service name %s. error: %s", serviceName, err)
					ticker.Reset(5 * time.Second)
					break
				}

				ok := updateEndPoint(serviceName, config, services, consulParse)

				ticker.Reset(time.Duration(config.SecondsToReloadServicesName) * time.Second)
				if ok {
					fmt.Printf("refresh service name %s successfully: %s\n", serviceName, time.Now().UTC())
					servicesNameDone <- ok
				} else {
					fmt.Printf("service name %s not found. Refresh was not successfully: %s\n", serviceName, time.Now().UTC())
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func updateEndPoint(
	serviceName string,
	config *config.Config,
	services []*consul.CatalogService,
	consulParse common_consul.Parse,
) bool {

	if len(services) == 0 {
		return false
	}

	service := services[rand.Intn(len(services))]

	address := "localhost"
	if config.Production {
		address = service.ServiceAddress
	}

	https := fmt.Sprintf("https://%s:%s", address, strconv.Itoa(service.ServicePort))
	fmt.Println("https selected: ", https)

	host := fmt.Sprintf("%s:%s", address, strconv.Itoa(service.ServicePort))
	fmt.Println("host selected: ", host)

	switch consulParse {
	case common_consul.CertificatesAndSecurityKeys:
		config.Certificates.EndPointGetCertificateCA = fmt.Sprintf("%s/%s", https, config.Certificates.APIPathCertificateCA)
		config.Certificates.EndPointGetCertificateHost = fmt.Sprintf("%s/%s", https, config.Certificates.APIPathCertificateHost)
		config.Certificates.EndPointGetCertificateHostKey = fmt.Sprintf("%s/%s", https, config.Certificates.APIPathCertificateHostKey)
		config.SecurityKeys.EndPointGetPublicKeys = fmt.Sprintf("%s/%s", https, config.SecurityKeys.APIPathPublicKeys)
		return true

	case common_consul.SecurityRSAKeys:
		config.SecurityRSAKeys.EndPointGetRSAPublicKeys = fmt.Sprintf("%s/%s", https, config.SecurityRSAKeys.APIPathRSAPublicKeys)
		return true

	case common_consul.EmailService:
		config.EmailService.Host = host
		return true

	default:
		return false
	}
}
