package main

import (
	"TOomaAh/emby_exporter_go/emby"
	"TOomaAh/emby_exporter_go/geoip"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Coordinate struct {
	latitude  string
	longitude string
}

var Options struct {
	Emby   string `short:"e" long:"emby" description:"emby url (don't forget scheme)" default:"http://localhost:8096"`
	Token  string `short:"t" long:"token" required:"true"`
	UserID string `short:"u" long:"user-id" required:"true"`
	Port   int    `short:"p" long:"port" default:"9210"`
	Sleep  int    `short:"s" long:"sleep" default:"10"`
}

var (
	embyInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "emby_system_info",
		Help: "All Emby Info",
	}, []string{"version", "wanAdress", "localAdress", "hasUpdateAvailable"})

	embyMediaItem = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "emby_media_item",
		Help: "All Media Item",
	}, []string{"name", "collectionType"})

	embySession = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "emby_sessions",
		Help: "All session",
	}, []string{"username", "client", "isPaused", "remoteEndPoint", "nowPlayingItemName", "nowPlayingItemType", "latitude", "longitude"})
)

func updateMetrics(emby emby.Emby, sleep time.Duration) {
	log.Printf("Server listening on 0.0.0.0:%d\n", Options.Port)

	for {
		info, _ := emby.GetSystemInfo()
		item, _ := emby.GetMediaItem()
		session, err := emby.GetSessions()

		if err != nil {
			fmt.Println(err)
		}

		embyInfo.With(prometheus.Labels{
			"version":            info.Version,
			"wanAdress":          info.WanAddress,
			"localAdress":        info.LocalAddress,
			"hasUpdateAvailable": strconv.FormatBool(info.HasUpdateAvailable)})

		for _, s := range *item {
			embyMediaItem.With(prometheus.Labels{
				"name":           s.Name,
				"collectionType": s.CollectionType,
			})
		}

		for _, s := range *session {
			if s.NowPlayingItem.Name != "" {

				geoip := geoip.New(s.RemoteEndPoint)
				geoInformation, _ := geoip.GetInfo()

				embySession.With(prometheus.Labels{
					"username":           s.UserName,
					"client":             s.Client,
					"remoteEndPoint":     s.RemoteEndPoint,
					"nowPlayingItemName": s.NowPlayingItem.Name,
					"nowPlayingItemType": s.NowPlayingItem.Type,
					"latitude":           strconv.FormatFloat(geoInformation.Lat, 'f', -1, 32),
					"longitude":          strconv.FormatFloat(geoInformation.Lon, 'f', -1, 32),
					"isPaused":           strconv.FormatBool(s.PlayState.IsPaused),
				})
			}
		}
		log.Println("metrics updated")
		time.Sleep(sleep * time.Second)
	}

}

func registerMetrics(r *prometheus.Registry) {
	// Metrics have to be registered to be exposed:
	r.MustRegister(embyInfo)
	r.MustRegister(embyMediaItem)
	r.MustRegister(embySession)
}

func main() {
	_, err := flags.ParseArgs(&Options, os.Args)

	if err != nil {
		log.Fatalln(err)
	}

	sleepDuration := time.Duration(Options.Sleep)

	var emby = emby.New(Options.Emby, Options.Token, Options.UserID)
	go updateMetrics(emby, sleepDuration)
	r := prometheus.NewRegistry()
	registerMetrics(r)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", Options.Port), nil)

}
