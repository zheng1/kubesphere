package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"
	"k8s.io/klog"
	"kubesphere.io/kubesphere/pkg/simple/client"
	"kubesphere.io/kubesphere/pkg/simple/client/prometheus"
	"kubesphere.io/kubesphere/pkg/utils/signals"
)

func main() {
	var whSvrParameters WhSvrParameters
	var aSvrParameters AccessSvrParameters

	prometheusOptions := prometheus.NewPrometheusOptions()
	prometheusOptions.AddFlags(flag.CommandLine)

	// get command line parameters
	flag.IntVar(&whSvrParameters.port, "port", 443, "Webhook server port.")
	flag.IntVar(&aSvrParameters.port, "accessPort", 80, "Access server port.")
	flag.StringVar(&whSvrParameters.certFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&whSvrParameters.keyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	csop := client.NewClientSetOptions()
	csop.SetPrometheusOptions(prometheusOptions)

	client.NewClientSetFactory(csop, signals.SetupSignalHandler())

	whsvr := startWebHookServer(whSvrParameters)
	asvr := startAccessServer(aSvrParameters)

	klog.Info("Server started")

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	klog.Infof("Got OS shutdown signal, shutting down webhook server gracefully...")
	_ = whsvr.server.Shutdown(context.Background())
	_ = asvr.server.Shutdown(context.Background())
}
