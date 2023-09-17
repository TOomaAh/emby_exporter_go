package emby

import (
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
	Url        string
	Token      string
	UserID     string
	Port       string
	GeoIp      bool
	httpClient *http.Client
}

type ServerMetrics struct {
	Info           *SystemInfo
	LibraryMetrics []*LibraryMetrics
	Sessions       []*SessionsMetrics
	SessionsCount  int
	Activity       []*ActivityMetric
	Alert          []*AlertMetrics
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

func (s *Server) GetServerInfo() *SystemInfo {
	var systemInfo SystemInfo
	err := s.request("GET", "/System/Info", "", &systemInfo)
	if err != nil {
		log.Println("Emby Server - GetServerInfo : " + err.Error())
		return nil
	}

	return &systemInfo
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
