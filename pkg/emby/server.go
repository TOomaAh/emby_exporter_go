package emby

import (
	"TOomaAh/emby_exporter_go/internal/entity"
	"TOomaAh/emby_exporter_go/pkg/logger"
	"TOomaAh/emby_exporter_go/pkg/request"
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
	client *request.Client
	UserID string
	Logger logger.Interface
}

type ServerInfo struct {
	Hostname string
	Port     string
	UserID   string
	Token    string
}

func NewServer(s *ServerInfo, logger logger.Interface) *Server {
	if s.Hostname == "" {
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

	client, err := request.NewClient(s.Hostname+":"+s.Port, s.Token)

	if err != nil {
		logger.Fatal("Cannot create a new client for Emby Server " + err.Error())
		os.Exit(1)
	}

	server := &Server{
		UserID: s.UserID,
		client: client,
		Logger: logger,
	}

	return server
}

func (s Server) GetSessions() (*[]entity.Sessions, error) {
	var sessions []entity.Sessions

	request, err := s.client.NewRequest(http.MethodGet, sessionPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.client.Do(request, &sessions)

	if err != nil {
		return nil, err
	}

	return &sessions, nil
}

func (s *Server) GetActivity() (*entity.Activity, error) {
	var activity entity.Activity
	request, err := s.client.NewRequest(http.MethodGet, activityPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.client.Do(request, &activity)

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

func (s *Server) Ping() error {
	request, err := s.client.NewRequest(http.MethodGet, pingPath, nil)

	if err != nil {
		return err
	}

	err = s.client.Do(request, nil)

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) GetLibrary() (*entity.LibraryInfo, error) {
	var library entity.LibraryInfo

	request, err := s.client.NewRequest(http.MethodGet, libraryPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.client.Do(request, &library)

	if err != nil {
		return nil, err
	}

	return &library, nil
}

func (s *Server) GetServerInfo() (*entity.SystemInfo, error) {
	var systemInfo entity.SystemInfo
	request, err := s.client.NewRequest(http.MethodGet, systemInfoPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.client.Do(request, &systemInfo)

	if err != nil {
		return nil, err
	}

	return &systemInfo, nil
}

func (s *Server) GetLibrarySize(itemID, contentType string) (int, error) {
	var librarySize int
	var library entity.Library

	request, err := s.client.NewRequest(http.MethodGet, "/Users/"+
		s.UserID+
		"/Items?IncludeItemTypes=Movie&Recursive=true&Fields=BasicSyncInfo&EnableImageTypes=Primary&ParentId="+
		itemID+"&Limit=1&IncludeItemTypes="+includeType[contentType], nil)

	if err != nil {
		return 0, err
	}

	err = s.client.Do(request, &library)

	if err != nil {
		return 0, err
	}

	librarySize = library.TotalRecordCount

	return librarySize, nil
}

func (s *Server) GetAlerts() (*entity.Alert, error) {
	var alert entity.Alert
	request, err := s.client.NewRequest(http.MethodGet, systemAlertPath, nil)

	if err != nil {
		return nil, err
	}

	err = s.client.Do(request, &alert)

	if err != nil {
		return nil, err
	}
	return &alert, nil
}
