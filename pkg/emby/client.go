package emby

import (
	"TOomaAh/emby_exporter_go/internal/entity"
	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
)

type EmbyClient struct {
	Server        *Server
	ServerMetrics *entity.ServerMetrics
	Logger        *logger.Logger
}

func NewEmbyClient(s *Server) *EmbyClient {
	return &EmbyClient{
		Server:        s,
		ServerMetrics: &entity.ServerMetrics{},
		Logger:        logger.New("info"),
	}
}

func (s *Server) GetSessionsMetrics() []*entity.SessionsMetrics {
	var sessions []entity.Sessions
	err := s.request("GET", "/Sessions", "", &sessions)
	if err != nil {
		s.Logger.Info("Cannot get sessions, maybe your server is unreachable " + err.Error())
		return []*entity.SessionsMetrics{}
	}

	count := 0
	for i := 0; i < len(sessions); i++ {
		if sessions[i].HasPlayMethod() {
			count++
		}
	}

	var sessionResult []*entity.SessionsMetrics = make([]*entity.SessionsMetrics, count)
	count = 0
	db := geoip.GetGeoIPDatabase()
	var sessionMetrics *entity.SessionsMetrics

	//To retrieve only the playback sessions and not the connected devices
	for _, session := range sessions {
		if session.HasPlayMethod() {

			if err != nil {
				s.Logger.Info("Emby Server - GetSessions : " + err.Error())
				return []*entity.SessionsMetrics{}
			}

			sessionMetrics = session.To()

			if s.GeoIp {
				sessionMetrics.Latitude, sessionMetrics.Longitude = db.GetLocation(session.RemoteEndPoint)
				sessionMetrics.City = db.GetCity(session.RemoteEndPoint)
				sessionMetrics.Region = db.GetRegion(session.RemoteEndPoint)
				sessionMetrics.CountryCode = db.GetCountryCode(session.RemoteEndPoint)
			}

			sessionResult[count] = sessionMetrics
			count++
		}
	}

	return sessionResult
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
