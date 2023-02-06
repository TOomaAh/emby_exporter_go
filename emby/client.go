package emby

import (
	"log"
)

type EmbyClient struct {
	Server *Server
}

func NewEmbyClient(s *Server) *EmbyClient {
	return &EmbyClient{
		Server: s,
	}
}

func (c *EmbyClient) GetMetrics() *ServerMetrics {
	serverMetrics := ServerMetrics{}
	systemInfo, err := c.Server.GetServerInfo()

	if err != nil {
		return nil
	}

	serverMetrics.Info = systemInfo

	library, err := c.Server.GetLibrary()

	if err != nil {
		return nil
	}

	libraryMetrics := []LibraryMetrics{}

	for _, l := range library.LibraryItem {
		size, _ := c.Server.GetLibrarySize(&l)
		serverMetrics.LibraryMetrics = append(libraryMetrics, LibraryMetrics{
			Name: l.Name,
			Size: size,
		})
	}

	sessions, err := c.Server.GetSessions()
	if err != nil {
		log.Println("Emby Client - GetMetrics : " + err.Error())
		return nil
	}
	serverMetrics.Sessions = *sessions
	serverMetrics.SessionsCount = len(*sessions)

	activity, err := c.Server.GetActivity()

	if err == nil {
		for _, a := range activity.Items {
			serverMetrics.Activity = append(serverMetrics.Activity, ActivityMetric{
				ID:       a.ID,
				Name:     a.Name,
				Type:     a.Type,
				Severity: a.Severity,
				Date:     a.Date,
			})
		}
	}

	alert, err := c.Server.GetAlert()

	if err == nil {
		for _, a := range alert.Items {
			serverMetrics.Alert = append(serverMetrics.Alert, AlertMetrics{
				ID:            a.ID,
				Name:          a.Name,
				Overview:      a.Overview,
				ShortOverview: a.ShortOverview,
				Type:          a.Type,
				Date:          a.Date,
				Severity:      a.Severity,
			})
		}
	}

	return &serverMetrics

}
