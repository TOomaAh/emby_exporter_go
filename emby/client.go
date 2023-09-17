package emby

import "fmt"

type EmbyClient struct {
	Server        *Server
	ServerMetrics *ServerMetrics
}

func NewEmbyClient(s *Server) *EmbyClient {
	return &EmbyClient{
		Server:        s,
		ServerMetrics: &ServerMetrics{},
	}
}

func (c *EmbyClient) GetMetrics() *ServerMetrics {
	var serverMetrics *ServerMetrics = c.ServerMetrics

	systemInfo := c.Server.GetServerInfo()

	if systemInfo == nil {
		systemInfo = &SystemInfo{
			Version:            "0.0.0",
			HasPendingRestart:  false,
			HasUpdateAvailable: false,
			LocalAddress:       "",
			WanAddress:         "",
		}
	}

	serverMetrics.Info = systemInfo

	libraries := c.Server.GetLibrary()
	fmt.Println("len(libraries.LibraryItem): ", len(libraries.LibraryItem), "len(serverMetrics.LibraryMetrics): ", len(serverMetrics.LibraryMetrics))
	if len(libraries.LibraryItem) != len(serverMetrics.LibraryMetrics) {
		serverMetrics.LibraryMetrics = make([]*LibraryMetrics, len(libraries.LibraryItem))
	}

	for i, l := range libraries.LibraryItem {
		serverMetrics.LibraryMetrics[i] = &LibraryMetrics{
			Name: l.Name,
			Size: c.Server.GetLibrarySize(&l),
		}
	}

	c.ServerMetrics.Sessions = c.Server.GetSessionsMetrics()
	c.ServerMetrics.SessionsCount = len(c.ServerMetrics.Sessions)

	activity := c.Server.GetActivity()

	if len(activity.Items) != len(serverMetrics.Activity) || len(serverMetrics.Activity) == 0 {
		serverMetrics.Activity = make([]*ActivityMetric, len(activity.Items))
	}

	for i, a := range activity.Items {
		serverMetrics.Activity[i] = &ActivityMetric{
			ID:       a.ID,
			Name:     a.Name,
			Type:     a.Type,
			Severity: a.Severity,
			Date:     a.Date,
		}
	}

	alert := c.Server.GetAlert()

	if len(alert.Items) != len(serverMetrics.Alert) || len(serverMetrics.Alert) == 0 {
		serverMetrics.Alert = make([]*AlertMetrics, len(alert.Items))
	}

	for i, a := range alert.Items {
		serverMetrics.Alert[i] = &AlertMetrics{
			ID:            a.ID,
			Overview:      a.Overview,
			ShortOverview: a.ShortOverview,
			Type:          a.Type,
			Date:          a.Date,
			Severity:      a.Severity,
		}
	}
	return nil
}
