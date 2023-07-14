package emby

import (
	"log"
)

type EmbyClient struct {
	Server        *Server
	ServerMetrics *ServerMetrics
}

func cleanMetrics(c *EmbyClient) {
	c.ServerMetrics.LibraryMetrics = nil
	c.ServerMetrics.Activity = nil
	c.ServerMetrics.Alert = nil
}

func NewEmbyClient(s *Server) *EmbyClient {
	return &EmbyClient{
		Server:        s,
		ServerMetrics: &ServerMetrics{},
	}
}

func (c *EmbyClient) GetMetrics() error {

	cleanMetrics(c)
	systemInfo, err := c.Server.GetServerInfo()

	if err != nil {
		return nil
	}

	c.ServerMetrics.Info = systemInfo

	library, err := c.Server.GetLibrary()

	if err != nil {
		return nil
	}

	for _, l := range library.LibraryItem {
		size, _ := c.Server.GetLibrarySize(&l)
		c.ServerMetrics.LibraryMetrics = append(c.ServerMetrics.LibraryMetrics, LibraryMetrics{
			Name: l.Name,
			Size: size,
		})
	}

	sessions, err := c.Server.GetSessions()
	if err != nil {
		log.Println("Emby Client - GetMetrics : " + err.Error())
		return nil
	}
	c.ServerMetrics.Sessions = *sessions
	c.ServerMetrics.SessionsCount = len(*sessions)

	activity, err := c.Server.GetActivity()

	if err == nil {
		for _, a := range activity.Items {
			c.ServerMetrics.Activity = append(c.ServerMetrics.Activity, ActivityMetric{
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
			c.ServerMetrics.Alert = append(c.ServerMetrics.Alert, AlertMetrics{
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
	return nil
}
