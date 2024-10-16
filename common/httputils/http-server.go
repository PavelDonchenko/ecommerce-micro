package httputils

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/PavelDonchenko/ecommerce-micro/common/certificates"
	"github.com/PavelDonchenko/ecommerce-micro/common/config"
)

type HttpServer interface {
	RunTLSServer(serverErr chan error) *http.Server
	RunUnsecuredServer(serverErr chan error) *http.Server
}

var srv *http.Server

type httpServer struct {
	config *config.Config
	router *gin.Engine
	cert   certificates.Certificator
}

func NewHttpServer(
	config *config.Config,
	router *gin.Engine,
	cert certificates.Certificator,
) *httpServer {
	return &httpServer{
		config: config,
		router: router,
		cert:   cert,
	}
}

func (s *httpServer) RunUnsecuredServer(serverErr chan error) *http.Server {
	srv = &http.Server{
		Addr:         s.config.ListenPort,
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	log.Printf("Uncecured server started, istening on port %s", s.config.ListenPort)

	return srv
}

func (s *httpServer) RunTLSServer(serverErr chan error) *http.Server {
	var err error
	srv = s.mountTLSServer()

	go func() {
		if err = srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	log.Printf("TLS server started, listening on port %s", s.config.ListenPort)

	return srv
}

func (s *httpServer) mountTLSServer() *http.Server {
	return &http.Server{
		Addr:         s.config.ListenPort,
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			GetCertificate: s.cert.GetLocalCertificate,
			ClientCAs:      s.cert.GetLocalCertificateCA(),
		},
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
}
