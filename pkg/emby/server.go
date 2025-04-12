package emby

import (
	"TOomaAh/emby_exporter_go/internal/entity"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"TOomaAh/emby_exporter_go/pkg/request"
	"errors"
	"net/http"
	"os"
)

var includeType = map[string]string{
	"movies":  "Movie",
	"tvshows": "Series",
	"boxsets": "BoxSet",
	"music":   "MusicArtist",
	"Episode": "TV Show",
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
	UserID string

	IsRecheable    bool
	requestManager *request.RequestManager

	Logger logger.Interface
}

type ServerInfo struct {
	Hostname string
	Port     string
	UserID   string
	Token    string
}

func NewServer(s *ServerInfo, logger logger.Interface) *Server {
	if s.Hostname == "" || s.Hostname == "http://" {
		s.Hostname = "http://localhost"
	}

	if s.Port == "" {
		s.Port = "8096"
	}

	if s.UserID == "" {
		logger.Fatal("UserID is not set in the configuration file")
		os.Exit(1)
	}

	if s.Token == "" {
		logger.Fatal("Token is not set in the configuration file")
		os.Exit(1)
	}

	client, err := NewEmbyClient(s.Hostname+":"+s.Port, s.Token)

	if err != nil {
		logger.Fatal("Cannot create a new client for Emby Server " + err.Error())
		os.Exit(1)
	}

	server := &Server{
		UserID:         s.UserID,
		IsRecheable:    false,
		Logger:         logger,
		requestManager: request.NewRequestManager(client),
	}

	return server
}

func (s Server) GetSessions() (*[]entity.Sessions, error) {

	var sessions []entity.Sessions

	req, err := s.requestManager.NewRequest(http.MethodGet, sessionPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.requestManager.Do(req, &sessions)

	if err != nil {

		if !errors.Is(err, request.ErrorCannotReadBody) {
			s.IsRecheable = false
		}

		return nil, err
	}

	return &sessions, nil
}

func (s *Server) GetActivity() (*entity.Activity, error) {
	var activity entity.Activity
	req, err := s.requestManager.NewRequest(http.MethodGet, activityPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.requestManager.Do(req, &activity)

	if err != nil {

		if !errors.Is(err, request.ErrorCannotReadBody) {
			s.IsRecheable = false
		}

		return nil, err
	}

	return &activity, nil
}

func (s *Server) Ping() error {
	req, err := s.requestManager.NewRequest(http.MethodGet, pingPath, nil)

	if err != nil {
		return err
	}

	err = s.requestManager.Do(req, nil)

	if err != nil {

		if !errors.Is(err, request.ErrorCannotReadBody) {
			s.IsRecheable = false
		}

		return err
	}

	return nil
}

func (s *Server) GetLibrary() (*entity.LibraryInfo, error) {
	var library entity.LibraryInfo

	req, err := s.requestManager.NewRequest(http.MethodGet, libraryPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.requestManager.Do(req, &library)

	if err != nil {

		if !errors.Is(err, request.ErrorCannotReadBody) {
			s.IsRecheable = false
		}

		return nil, err
	}

	return &library, nil
}

func (s *Server) GetServerInfo() (*entity.SystemInfo, error) {
	var systemInfo entity.SystemInfo
	req, err := s.requestManager.NewRequest(http.MethodGet, systemInfoPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.requestManager.Do(req, &systemInfo)

	if err != nil {

		if !errors.Is(err, request.ErrorCannotReadBody) {
			s.IsRecheable = false
		}

		return nil, err
	}

	return &systemInfo, nil
}

func (s *Server) GetLibrarySize(itemID, contentType string) (int, error) {
	var librarySize int
	var library entity.Library

	req, err := s.requestManager.NewRequest(http.MethodGet, "/Users/"+
		s.UserID+
		"/Items?IncludeItemTypes=Movie&Recursive=true&Fields=BasicSyncInfo&EnableImageTypes=Primary&ParentId="+
		itemID+"&Limit=1&IncludeItemTypes="+includeType[contentType], nil)

	if err != nil {
		return 0, err
	}

	err = s.requestManager.Do(req, &library)

	if err != nil {

		if !errors.Is(err, request.ErrorCannotReadBody) {
			s.IsRecheable = false
		}

		return 0, err
	}

	librarySize = library.TotalRecordCount

	return librarySize, nil
}

func (s *Server) GetAlerts() (*entity.Alert, error) {
	var alert entity.Alert
	req, err := s.requestManager.NewRequest(http.MethodGet, systemAlertPath, nil)

	if err != nil {

		return nil, err
	}

	err = s.requestManager.Do(req, &alert)

	if err != nil {

		if !errors.Is(err, request.ErrorCannotReadBody) {
			s.IsRecheable = false
		}

		return nil, err
	}
	return &alert, nil
}
