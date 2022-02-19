package emby

import "log"

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
		log.Fatal(err)
		return nil
	}

	serverMetrics.Info = systemInfo

	library, err := c.Server.GetLibrary()

	if err != nil {
		log.Fatal("error pour les libs", err)
		return nil
	}

	libraryMetrics := []LibraryMetrics{}

	for _, l := range *&library.LibraryItem {
		size, _ := c.Server.GetLibrarySize(&l)
		libraryMetrics = append(libraryMetrics, LibraryMetrics{
			Name: l.Name,
			Size: size,
		})
	}

	serverMetrics.LibraryMetrics = libraryMetrics

	sessions, err := c.Server.GetSessions()
	if err != nil {
		log.Fatal("error pour les sessions", err)
		return nil
	}
	serverMetrics.Sessions = *sessions
	serverMetrics.SessionsCount = len(*sessions)

	return &serverMetrics

}
