package main

import (
	"TOomaAh/emby_exporter_go/emby"
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

var Options struct {
	Emby   string `short:"e" long:"emby" description:"emby url (don't forget scheme)" default:"http://localhost:8096"`
	Token  string `short:"t" long:"token" required:"true"`
	UserID string `short:"u" long:"user-id" required:"true"`
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
	}, []string{"username", "client", "remoteEndPoint", "nowPlayingItemName", "nowPlayingItemType"})
)

func updateMetrics(emby emby.Emby, sleep time.Duration) {

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
				embySession.With(prometheus.Labels{
					"username":           s.UserName,
					"client":             s.Client,
					"remoteEndPoint":     s.RemoteEndPoint,
					"nowPlayingItemName": s.NowPlayingItem.Name,
					"nowPlayingItemType": s.NowPlayingItem.Type,
				})
			}
		}
		log.Println("metrics updated")
		time.Sleep(sleep * time.Second)
	}

}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.NewRegistry()
	prometheus.MustRegister(embyInfo)
	prometheus.MustRegister(embyMediaItem)
	prometheus.MustRegister(embySession)
}

func main() {

	_, err := flags.ParseArgs(&Options, os.Args)

	if err != nil {
		log.Fatalln(err)
	}

	sleepDuration := time.Duration(Options.Sleep)

	var emby = emby.New(Options.Emby, Options.Token, Options.UserID)

	go updateMetrics(emby, sleepDuration)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}
