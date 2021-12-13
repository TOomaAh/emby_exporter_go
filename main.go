package main

import (
	"TOomaAh/emby_exporter_go/emby"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

func updateMetrics(emby emby.Emby) {

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

}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.NewRegistry()
	prometheus.MustRegister(embyInfo)
	prometheus.MustRegister(embyMediaItem)
	prometheus.MustRegister(embySession)
}

func main() {
	var emby = emby.New("", "", "")
	updateMetrics(emby)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}
