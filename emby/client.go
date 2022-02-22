package emby

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
		libraryMetrics = append(libraryMetrics, LibraryMetrics{
			Name: l.Name,
			Size: size,
		})
	}

	serverMetrics.LibraryMetrics = libraryMetrics

	sessions, err := c.Server.GetSessions()
	if err != nil {
		return nil
	}
	serverMetrics.Sessions = *sessions
	serverMetrics.SessionsCount = len(*sessions)

	alert, err := c.Server.GetAlert()

	if err == nil {
		for _, a := range alert.Items {
			serverMetrics.Activity = append(serverMetrics.Activity, ActivityMetric{
				Name:     a.Name,
				Type:     a.Type,
				Severity: a.Severity,
				Date:     a.Date,
			})
		}
	}

	return &serverMetrics

}
