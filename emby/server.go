package emby

import (
	"TOomaAh/emby_exporter_go/geoip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var includeType = map[string]string{
	"movies":  "Movie",
	"tvshows": "Series",
	"boxsets": "BoxSet",
	"music":   "MusicArtist",
	"Episode": "TV Show",
}

type Server struct {
	Url    string
	Token  string
	UserID string
	Port   int
	GeoIp  bool
}

func NewServer(url, token, userID string, port int, geoip bool) *Server {

	server := &Server{
		Url:    url,
		Token:  token,
		UserID: userID,
		Port:   port,
		GeoIp:  geoip,
	}

	return server
}

func (s *Server) GetServerInfo() (*SystemInfo, error) {
	resp, err := s.request("GET", "/System/Info", "")
	if err != nil {
		log.Println("Emby Server - GetServerInfo : " + err.Error())
		return nil, err
	}

	var systemInfo SystemInfo
	err = json.Unmarshal(resp, &systemInfo)
	if err != nil {
		log.Println("GetServerInfo : " + err.Error())
		return nil, err
	}

	return &systemInfo, nil
}

func (s *Server) GetLibrary() (*LibraryInfo, error) {
	resp, err := s.request("GET", "/Library/VirtualFolders/Query", "")
	if err != nil {
		return nil, err
	}
	var library LibraryInfo
	err = json.Unmarshal(resp, &library)
	if err != nil {
		log.Println("Emby Server - GetLibrary : " + err.Error())
		return nil, err
	}

	return &library, nil
}

func (s *Server) GetSessions() (*[]SessionsMetrics, error) {
	resp, err := s.request("GET", "/Sessions", "")
	if err != nil {
		log.Println("Cannot get sessions, maybe your server is unreachable")
		return nil, err
	}

	var sessions []Sessions
	err = json.Unmarshal(resp, &sessions)
	if err != nil {
		log.Println("Emby Server - GetSessions : " + err.Error())
		return nil, err
	}
	var sessionResult []SessionsMetrics

	//To retrieve only the playback sessions and not the connected devices
	for _, session := range sessions {
		if session.PlayState.PlayMethod != "" {

			ip := geoip.New(session.RemoteEndPoint)
			var lat, long float64 = 0.0, 0.0
			var playbackPercent int64

			var city, region, countryCode, tvShow, season string

			var positionTicksSeconds float64 = (float64(session.PlayState.PositionTicks) / 10000000) / 60
			var totalTime float64 = (float64(session.NowPlayingItem.RunTimeTicks) / 10000000) / 60
			playbackPercent = (int64((positionTicksSeconds * 100) / totalTime))

			if session.NowPlayingItem.Type == "Episode" {
				tvShow = session.NowPlayingItem.SeriesName
				season = session.NowPlayingItem.SeasonName
			}

			if err == nil {
				if s.GeoIp {
					information, errGeoIp := ip.GetInfo()
					if errGeoIp == nil {
						city = information.City
						region = information.RegionName
						countryCode = information.CountryCode
						lat = information.Lat
						long = information.Lon
					}
				}
			}
			sessionResult = append(sessionResult, SessionsMetrics{
				Username:           session.UserName,
				Client:             session.Client,
				IsPaused:           session.PlayState.IsPaused,
				RemoteEndPoint:     session.RemoteEndPoint,
				Latitude:           lat,
				Longitude:          long,
				City:               city,
				Region:             region,
				CountryCode:        countryCode,
				NowPlayingItemName: session.NowPlayingItem.Name,
				NowPlayingItemType: session.NowPlayingItem.Type,
				MediaDuration:      session.NowPlayingItem.RunTimeTicks,
				PlaybackPosition:   session.PlayState.PositionTicks,
				PlaybackPercent:    playbackPercent,
				TVShow:             tvShow,
				Season:             season,
				PlayMethod:         session.PlayState.PlayMethod,
			})
		}
	}

	return &sessionResult, nil
}

func (s *Server) GetAlert() (*Alert, error) {
	resp, err := s.request("GET", "/System/ActivityLog/Entries?StartIndex=0&Limit=4&hasUserId=false", "")
	if err != nil {
		log.Println("Cannot get alert, maybe your server is unreachable")
		return nil, err
	}

	var alert Alert
	err = json.Unmarshal(resp, &alert)

	if err != nil {
		log.Println("Cannot parse alert response, your token is probably wrong")
		return nil, err
	}

	return &alert, nil
}

func (s *Server) GetActivity() (*Activity, error) {
	resp, err := s.request("GET", "/System/ActivityLog/Entries?StartIndex=0&Limit=7", "")
	if err != nil {
		log.Println("Cannot get activity, maybe your server is unreachable")
		return nil, err
	}

	var activity Activity
	err = json.Unmarshal(resp, &activity)

	if err != nil {
		log.Println("Cannot parse activity response, your token is probably wrong")
		return nil, err
	}

	return &activity, nil
}

func (s *Server) GetSessionsSize() (int, error) {
	sessions, err := s.GetSessions()
	if err != nil {
		log.Println("Cannot get session size, maybe your server is unreachable")
		return 0, err
	}

	return len(*sessions), nil
}

func (s *Server) GetLibrarySize(libraryItem *LibraryItem) (int, error) {
	var librarySize int
	resp, err := s.request("GET", fmt.Sprintf(
		//Ok I need minimum information. Only one Item and api returns the total number of items
		"/Users/%s/Items?IncludeItemTypes=Movie&Recursive=true&Fields=BasicSyncInfo&EnableImageTypes=Primary&ParentId=%s&Limit=1&IncludeItemTypes=%s",
		s.UserID,
		libraryItem.ItemID,
		includeType[libraryItem.LibraryOptions.ContentType]), "")

	if err != nil {
		log.Println("Cannot get library size, maybe your server is unreachable or your user is not allowed to access this library")
		return 0, err
	}

	var library Library
	err = json.Unmarshal(resp, &library)

	if err != nil {
		log.Println("Cannot parse library size response, your user id is probably wrong")
		return 0, err
	}
	librarySize = library.TotalRecordCount

	return librarySize, nil
}

func (s *Server) Ping() error {
	_, err := s.request("GET", "/System/Ping", "")
	return err
}

func (s *Server) request(method string, path string, body string) ([]byte, error) {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s:%d%s", s.Url, s.Port, path), strings.NewReader(body))
	req.Header.Set("X-Emby-Token", s.Token)
	req.Header.Set("Application-Type", "application/json")

	if len(body) > 0 {
		bodybytes := []byte(body)
		buf := bytes.NewBuffer(bodybytes)
		req.Body = ioutil.NopCloser(buf)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Problem with request to Emby Server")
		return nil, err
	}
	defer resp.Body.Close()
	respbody, _ := ioutil.ReadAll(resp.Body)
	return respbody, nil
}
