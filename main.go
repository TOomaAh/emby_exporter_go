package main

import (
	"TOomaAh/emby_exporter_go/emby"
	"TOomaAh/emby_exporter_go/metrics"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Coordinate struct {
	latitude  string
	longitude string
}

var Options struct {
	Emby     string `short:"e" long:"emby" description:"emby url without port (don't forget scheme)" required:"true"`
	EmbyPort int    `short:"p" long:"embyport" description:"emby port" default:"8096" required:"true"`
	Token    string `short:"t" long:"token" required:"true"`
	UserID   string `short:"u" long:"user-id" required:"true"`
	Port     int    `short:"P" long:"port" default:"9210"`
}

func main() {
	_, err := flags.ParseArgs(&Options, os.Args)
	logger := log.Default()
	if err != nil {
		log.Fatalln(err)
	}

	embyServer := emby.NewServer(Options.Emby, Options.Token, Options.UserID, Options.EmbyPort)
	client := emby.NewEmbyClient(embyServer)
	embyCollector := metrics.NewEmbyCollector(client)
	newRegistry := prometheus.NewRegistry()
	newRegistry.MustRegister(embyCollector)
	handler := promhttp.HandlerFor(newRegistry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	logger.Printf("Beginning to serve on port %d", Options.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", Options.Port), nil)

}
