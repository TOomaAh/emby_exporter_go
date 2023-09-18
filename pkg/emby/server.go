package emby

import (
	"TOomaAh/emby_exporter_go/internal/entity"
	"bytes"
	"encoding/json"
	"io"
	"log"
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
}

func NewServer(url, token, userID string, port int, geoip bool) *Server {
	server := &Server{
		Url:    url,
		Token:  token,
		UserID: userID,
		Port:   strconv.Itoa(port),
		GeoIp:  geoip,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	return server
}

func (s *Server) GetActivity() *entity.Activity {
	var activity entity.Activity
	err := s.request("GET", "/System/ActivityLog/Entries?StartIndex=0&Limit=7", "", &activity)

	if err != nil {
		log.Println("Cannot get activity, maybe your server is unreachable")
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
		log.Println("Problem with request to Emby Server")
		return err
	}

	defer resp.Body.Close()

	if v == nil {
		return nil
	}

	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		log.Println("Cannot parse response from Emby Server")
		return err
	}

	return nil
}

func (s *Server) GetLibrary() *entity.LibraryInfo {
	var library entity.LibraryInfo
	err := s.request("GET", "/Library/VirtualFolders/Query", "", &library)

	if err != nil {
		log.Println("Emby Server - GetLibrary : " + err.Error())
		library.LibraryItem = []entity.LibraryItem{}
		return &library
	}

	return &library
}

func (s *Server) GetServerInfo() *entity.SystemInfo {
	var systemInfo entity.SystemInfo
	err := s.request("GET", "/System/Info", "", &systemInfo)
	if err != nil {
		log.Println("Emby Server - GetServerInfo : " + err.Error())
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
		log.Println("Cannot get library size, maybe your server is unreachable or your user is not allowed to access this library : " + err.Error())
		return 0
	}

	librarySize = library.TotalRecordCount

	return librarySize
}

func (s *Server) GetAlert() *entity.Alert {
	var alert entity.Alert
	err := s.request("GET", "/System/ActivityLog/Entries?StartIndex=0&Limit=4&hasUserId=false", "", &alert)

	if err != nil {
		log.Println("Cannot get alert, maybe your server is unreachable")
		alert.Items = make([]entity.AlertItem, 0)
		return &alert
	}

	return &alert
}
