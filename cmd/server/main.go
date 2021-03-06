package main

import (
	"flag"
	"github.com/mafr/k8s-admission-webhook/pkg/server"
	"github.com/mafr/k8s-admission-webhook/pkg/validator"
	log "github.com/sirupsen/logrus"
)

func main() {
	listenAddress := flag.String("l", ":8443", "the address to listen on")
	certFile := flag.String("c", "/etc/webhook/certs/cert.pem", "server certificate in PEM format")
	keyFile := flag.String("k", "/etc/webhook/certs/key.pem", "server private key in PEM format")
	plainHttp := flag.Bool("p", false, "serve on plain HTTP for testing")
	logLevel := flag.String("L", "info", "lowest level to write to the log")

	flag.Parse()

	server.MustInitLogger(*logLevel)

	log.Infof("listening on %s", *listenAddress)
	if *plainHttp {
		log.Warn("running in plain HTTP mode (will NOT work in Kubernetes!)")
	}

	val := validator.ValidatorConfig{}
	val.Add(validator.NewCpuValidator())
	val.Add(validator.NewMemValidator())
	val.Add(validator.NewReplicasValidator())

	httpServer := server.NewServer(*listenAddress, val)

	if *plainHttp {
		// This is for testing only, Kubernetes won't accept plain HTTP webhooks.
		log.Fatal(httpServer.ListenAndServe())
	} else {
		log.Fatal(httpServer.ListenAndServeTLS(*certFile, *keyFile))
	}
}
