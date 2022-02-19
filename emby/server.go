package emby

import (
	"TOomaAh/emby_exporter_go/geoip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Server struct {
	Url    string
	Token  string
	UserID string
	Port   int
}

func NewServer(url, token, userID string, port int) *Server {

	server := &Server{
		Url:    url,
		Token:  token,
		UserID: userID,
		Port:   port,
	}

	return server
}

func (s *Server) GetServerInfo() (*SystemInfo, error) {
	resp, err := s.request("GET", "/System/Info", "")
	if err != nil {
		return nil, err
	}

	var systemInfo SystemInfo
	err = json.Unmarshal(resp, &systemInfo)
	if err != nil {
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
		return nil, err
	}

	return &library, nil
}

func (s *Server) GetSessions() (*[]SessionsMetrics, error) {
	resp, err := s.request("GET", "/Sessions", "")
	if err != nil {
		return nil, err
	}

	var sessions []Sessions
	err = json.Unmarshal(resp, &sessions)
	if err != nil {
		return nil, err
	}
	var sessionResult []SessionsMetrics

	//To retrieve only the playback sessions and not the connected devices
	for _, session := range sessions {
		if session.PlayState.PlayMethod != "" {
			ip := geoip.New(session.RemoteEndPoint)
			information, err := ip.GetInfo()
			var lat, long float64 = 0.0, 0.0
			if err == nil {
				lat = information.Lat
				long = information.Lon
			}
			sessionResult = append(sessionResult, SessionsMetrics{
				Username:           session.UserName,
				Client:             session.Client,
				IsPaused:           session.PlayState.IsPaused,
				RemoteEndPoint:     session.RemoteEndPoint,
				Latitude:           lat,
				Longitude:          long,
				NowPlayingItemName: session.NowPlayingItem.Name,
				NowPlayingItemType: session.NowPlayingItem.Type,
			})
		}
	}

	return &sessionResult, nil
}

func (s *Server) GetSessionsSize() (int, error) {
	sessions, err := s.GetSessions()
	if err != nil {
		return 0, err
	}

	return len(*sessions), nil
}

func (s *Server) GetLibrarySize(libraryItem *LibraryItem) (int, error) {
	var librarySize int
	resp, err := s.request("GET", fmt.Sprintf(
		//Ok I need minimum information. Only one Item and api returns the total number of items
		"/Users/%s/Items?IncludeItemTypes=Movie&Recursive=true&Fields=BasicSyncInfo&EnableImageTypes=Primary&ParentId=%s&Limit=1",
		s.UserID,
		libraryItem.ItemID), "")

	if err != nil {
		return 0, err
	}

	var library Library
	err = json.Unmarshal(resp, &library)

	if err != nil {
		return 0, err
	}
	librarySize = library.TotalRecordCount

	return librarySize, nil
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
		return nil, err
	}
	defer resp.Body.Close()
	respbody, _ := ioutil.ReadAll(resp.Body)
	return respbody, nil
}
