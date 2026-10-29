package emby

import (
	"TOomaAh/emby_exporter_go/internal/entity"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"TOomaAh/emby_exporter_go/pkg/request"
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

var includeType = map[string]string{
	"movies":      "Movie",
	"tvshows":     "Series",
	"boxsets":     "BoxSet",
	"music":       "MusicArtist",
	"musicvideos": "MusicVideo",
	"homevideos":  "Video,Photo",
	"photos":      "Photo",
	"books":       "Book",
	"mixed":       "Movie,Series",
	"playlists":   "Playlist",
}

const (
	pingPath        = "/System/Ping"
	sessionPath     = "/Sessions?IncludeAllSessionsIfAdmin=true&IsPlaying=true"
	activityPath    = "/System/ActivityLog/Entries?StartIndex=0&Limit=7"
	libraryPath     = "/Library/VirtualFolders/Query"
	systemInfoPath  = "/System/Info"
	systemAlertPath = "/System/ActivityLog/Entries?StartIndex=0&Limit=4&hasUserId=false"
)

type Server struct {
	reachable      atomic.Bool
	pinging        atomic.Bool
	requestManager *request.RequestManager

	Logger logger.Interface
}

// IsReachable reports whether the last request to the server succeeded.
func (s *Server) IsReachable() bool {
	return s.reachable.Load()
}

// StartPingLoop pings the server every interval until it answers.
// Only one loop runs at a time: extra calls return immediately.
func (s *Server) StartPingLoop(interval time.Duration) {
	if !s.pinging.CompareAndSwap(false, true) {
		return
	}
	defer s.pinging.Store(false)

	for {
		if err := s.Ping(); err == nil {
			s.Logger.Info("Server is reachable")
			return
		}
		s.Logger.Error("Server is not reachable, retrying in %s...", interval)
		time.Sleep(interval)
	}
}

type ServerInfo struct {
	Hostname string
	Port     string
	Token    string
}

func NewServer(s *ServerInfo, logger logger.Interface) *Server {
	if s.Hostname == "http://" {
		s.Hostname = "http://localhost"
	}

	if s.Port == "" {
		s.Port = "8096"
	}

	if s.Token == "" {
		logger.Fatal("Token is not set in the configuration file")
	}

	client, err := NewEmbyClient(s.Hostname+":"+s.Port, s.Token)

	if err != nil {
		logger.Fatal("Cannot create a new client for Emby Server " + err.Error())
	}

	server := &Server{
		Logger:         logger,
		requestManager: request.NewRequestManager(client),
	}

	return server
}

func (s *Server) GetSessions() ([]entity.Sessions, error) {

	sessions, err := get[[]entity.Sessions](s, sessionPath)
	if err != nil {
		return nil, err
	}
	return *sessions, nil
}

func (s *Server) GetActivity() (*entity.Activity, error) {
	return get[entity.Activity](s, activityPath)
}

func (s *Server) Ping() error {
	req, err := s.requestManager.NewRequest(http.MethodGet, pingPath, nil)

	if err != nil {
		return err
	}

	err = s.requestManager.Do(req, nil)

	s.reachable.Store(err == nil)
	return err
}

func (s *Server) GetLibrary() (*entity.LibraryInfo, error) {
	return get[entity.LibraryInfo](s, libraryPath)
}

func (s *Server) GetServerInfo() (*entity.SystemInfo, error) {
	return get[entity.SystemInfo](s, systemInfoPath)
}

func (s *Server) GetLibrarySize(itemID, contentType string) (int, error) {
	q := url.Values{}
	q.Set("Recursive", "true")
	q.Set("ParentId", itemID)
	q.Set("Limit", "0")

	if t, ok := includeType[contentType]; ok {
		q.Set("IncludeItemTypes", t)
	} else {
		q.Set("IsFolder", "false")
	}

	library, err := get[entity.Library](s, "/Items?"+q.Encode())
	if err != nil {
		return 0, err
	}
	return library.TotalRecordCount, nil
}

func (s *Server) GetAlerts() (*entity.Alert, error) {
	return get[entity.Alert](s, systemAlertPath)
}

func get[T any](s *Server, path string) (*T, error) {
	req, err := s.requestManager.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var v T
	if err := s.requestManager.Do(req, &v); err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			s.reachable.Store(false)
		}
		return nil, err
	}

	return &v, nil
}
