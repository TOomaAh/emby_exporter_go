package emby

import (
	"TOomaAh/emby_exporter_go/internal/entity"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var includeType = map[string]string{
	"movies":  "Movie",
	"tvshows": "Series",
	"boxsets": "BoxSet",
	"music":   "MusicArtist",
	"Episode": "TV Show",
}

type Server struct {
	httpClient *http.Client
	Url        string
	Token      string
	UserID     string
	Port       string
	GeoIp      bool
	Logger     logger.Interface
}

func NewServer(url, token, userID string, port int, geoip bool, logger logger.Interface) *Server {
	server := &Server{
		Url:    url,
		Token:  token,
		UserID: userID,
		Port:   strconv.Itoa(port),
		GeoIp:  geoip,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Logger: logger,
	}

	return server
}

func (s Server) GetSessionsMetrics() []*entity.SessionsMetrics {
	var sessions []entity.Sessions
	err := s.request("GET", "/Sessions?IncludeAllSessionsIfAdmin=true&IsPlaying=true", "", &sessions)
	if err != nil {
		s.Logger.Info("Cannot get sessions, maybe your server is unreachable " + err.Error())
		return []*entity.SessionsMetrics{}
	}

	count := 0
	for i := 0; i < len(sessions); i++ {
		if sessions[i].HasPlayMethod() {
			count++
		}
	}

	var sessionResult []*entity.SessionsMetrics = make([]*entity.SessionsMetrics, count)
	count = 0
	db := geoip.GetGeoIPDatabase()
	var sessionMetrics *entity.SessionsMetrics

	//To retrieve only the playback sessions and not the connected devices
	for _, session := range sessions {
		if session.HasPlayMethod() {

			if err != nil {
				s.Logger.Info("Emby Server - GetSessions : " + err.Error())
				return []*entity.SessionsMetrics{}
			}

			sessionMetrics = session.To()

			if s.GeoIp {
				sessionMetrics.Latitude, sessionMetrics.Longitude = db.GetLocation(session.RemoteEndPoint)
				sessionMetrics.City = db.GetCity(session.RemoteEndPoint)
				sessionMetrics.Region = db.GetRegion(session.RemoteEndPoint)
				sessionMetrics.CountryCode = db.GetCountryCode(session.RemoteEndPoint)
			}

			sessionResult[count] = sessionMetrics
			count++
		}
	}

	return sessionResult
}

func (s *Server) GetActivity() *entity.Activity {
	var activity entity.Activity
	err := s.request("GET", "/System/ActivityLog/Entries?StartIndex=0&Limit=7", "", &activity)

	if err != nil {
		s.Logger.Info("Cannot get activity, maybe your server is unreachable")
		activity.Items = make([]entity.ActivityItem, 0)
		return &activity
	}

	return &activity
}

func (s *Server) Ping() error {
	return s.request("GET", "/System/Ping", "", nil)
}

func (s *Server) request(method string, path string, body string, v interface{}) error {
	req, _ := http.NewRequest(method, s.Url+":"+s.Port+path, strings.NewReader(body))
	req.Header.Set("X-Emby-Token", s.Token)
	req.Header.Set("Application-Type", "application/json")

	if len(body) > 0 {
		bodybytes := []byte(body)
		buf := bytes.NewBuffer(bodybytes)
		req.Body = io.NopCloser(buf)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		s.Logger.Info("Problem with request to Emby Server")
		return err
	}

	defer resp.Body.Close()

	if v == nil {
		return nil
	}

	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		s.Logger.Info("Cannot parse response from Emby Server")
		return err
	}

	return nil
}

func (s *Server) GetLibrary() *entity.LibraryInfo {
	var library entity.LibraryInfo
	err := s.request("GET", "/Library/VirtualFolders/Query", "", &library)

	if err != nil {
		s.Logger.Info("Emby Server - GetLibrary : " + err.Error())
		library.LibraryItem = []entity.LibraryItem{}
		return &library
	}

	return &library
}

func (s *Server) GetServerInfo() *entity.SystemInfo {
	var systemInfo entity.SystemInfo
	err := s.request("GET", "/System/Info", "", &systemInfo)
	if err != nil {
		s.Logger.Info("Emby Server - GetServerInfo : " + err.Error())
		return &entity.SystemInfo{
			Version:            "0.0.0",
			HasPendingRestart:  false,
			HasUpdateAvailable: false,
			LocalAddress:       "",
			WanAddress:         "",
		}
	}
	return &systemInfo
}

func (s *Server) GetLibrarySize(libraryItem *entity.LibraryItem) int {
	var librarySize int
	var library entity.Library
	err := s.request("GET",
		//Ok I need minimum information. Only one Item and api returns the total number of items
		"/Users/"+
			s.UserID+
			"/Items?IncludeItemTypes=Movie&Recursive=true&Fields=BasicSyncInfo&EnableImageTypes=Primary&ParentId="+
			libraryItem.ItemID+"&Limit=1&IncludeItemTypes="+includeType[libraryItem.LibraryOptions.ContentType], "", &library)

	if err != nil {
		s.Logger.Info("Cannot get library size, maybe your server is unreachable or your user is not allowed to access this library : " + err.Error())
		return 0
	}

	librarySize = library.TotalRecordCount

	return librarySize
}

func (s *Server) GetAlert() *entity.Alert {
	var alert entity.Alert
	err := s.request("GET", "/System/ActivityLog/Entries?StartIndex=0&Limit=4&hasUserId=false", "", &alert)

	if err != nil {
		s.Logger.Error("Cannot get alert, maybe your server is unreachable")
		alert.Items = make([]entity.AlertItem, 0)
		return &alert
	}

	return &alert
}
