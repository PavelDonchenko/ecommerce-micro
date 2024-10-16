package security

import (
	"crypto/x509"
	"errors"
	"time"

	"github.com/PavelDonchenko/ecommerce-micro/auth/service"
	common_cert "github.com/PavelDonchenko/ecommerce-micro/common/certificates"
	"github.com/PavelDonchenko/ecommerce-micro/common/config"
	"github.com/PavelDonchenko/ecommerce-micro/common/helpers"
)

type ManagerCertificates struct {
	config         *config.Config
	service        service.Certificator
	common_service common_cert.Certificator
}

var (
	certPath string
	keyPath  string
)

func NewManagerCertificates(
	config *config.Config,
	service service.Certificator,
	common_service common_cert.Certificator,
) *ManagerCertificates {
	certPath, keyPath = common_service.GetPathsCertificateHostAndKey()
	return &ManagerCertificates{
		config:         config,
		service:        service,
		common_service: common_service,
	}
}

func (m *ManagerCertificates) VerifyCertificates() bool {
	if !helpers.FileExists(certPath) || !helpers.FileExists(keyPath) {
		err := m.newCertificate()

		return err == nil
	}

	cert, err := m.common_service.ReadCertificate()
	if err != nil {
		return false
	}

	if cert == nil || cert.NotAfter.AddDate(0, 0, -7).Before(time.Now().UTC()) {
		err := m.newCertificate()
		if err != nil {
			return false
		}

	}

	return true
}

func (m *ManagerCertificates) GetCertificateCA() error {
	return errors.New("not implemented")
}

func (m *ManagerCertificates) GetCertificate() error {
	return errors.New("not implemented")
}

func (m *ManagerCertificates) newCertificate() error {
	caBytes, caPrivateKey, err := m.service.GenerateCertificateAuthority()
	if err != nil {
		return err
	}

	err = m.service.CreateCertificateCAPEM(caBytes)
	if err != nil {
		return err
	}

	err = m.service.CreateCertificateCAPrivateKeyPEM(caPrivateKey)
	if err != nil {
		return err
	}

	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return err
	}

	hostBytes, hostPrivateKey, err := m.service.GenerateCertificateHost(caCert, caPrivateKey)
	if err != nil {
		return err
	}

	err = m.service.CreateCertificateHostPEM(hostBytes)
	if err != nil {
		return err
	}

	err = m.service.CreateCertificateHostPrivateKeyPEM(hostPrivateKey)
	if err != nil {
		return err
	}

	return nil
}
