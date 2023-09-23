package emby

import (
	"TOomaAh/emby_exporter_go/internal/entity"
	"TOomaAh/emby_exporter_go/pkg/logger"
)

type EmbyClient struct {
	Server        *Server
	ServerMetrics *entity.ServerMetrics
	Logger        logger.Interface
}

func NewEmbyClient(s *Server, l logger.Interface) *EmbyClient {
	return &EmbyClient{
		Server:        s,
		ServerMetrics: &entity.ServerMetrics{},
		Logger:        l,
	}
}

func (c *EmbyClient) GetMetrics() *entity.ServerMetrics {
	var serverMetrics *entity.ServerMetrics = c.ServerMetrics

	serverMetrics.Info = c.Server.GetServerInfo()

	libraries := c.Server.GetLibrary()
	if len(libraries.LibraryItem) != len(serverMetrics.LibraryMetrics) {
		serverMetrics.LibraryMetrics = make([]*entity.LibraryMetrics, len(libraries.LibraryItem))
	}

	for i, l := range libraries.LibraryItem {
		serverMetrics.LibraryMetrics[i] = &entity.LibraryMetrics{
			Name: l.Name,
			Size: c.Server.GetLibrarySize(&l),
		}
	}

	c.ServerMetrics.Sessions = c.Server.GetSessionsMetrics()
	c.ServerMetrics.SessionsCount = len(c.ServerMetrics.Sessions)

	activity := c.Server.GetActivity()

	if len(activity.Items) != len(serverMetrics.Activity) || (len(serverMetrics.Activity) == 0 && len(activity.Items) > 0) {
		serverMetrics.Activity = make([]*entity.ActivityMetric, len(activity.Items))
	}

	for i, a := range activity.Items {
		serverMetrics.Activity[i] = &entity.ActivityMetric{
			ID:       a.ID,
			Name:     a.Name,
			Type:     a.Type,
			Severity: a.Severity,
			Date:     a.Date,
		}
	}

	alert := c.Server.GetAlert()

	if len(alert.Items) != len(serverMetrics.Alert) || (len(serverMetrics.Alert) == 0 && len(alert.Items) > 0) {
		serverMetrics.Alert = make([]*entity.AlertMetrics, len(alert.Items))
	}

	for i, a := range alert.Items {
		serverMetrics.Alert[i] = &entity.AlertMetrics{
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
