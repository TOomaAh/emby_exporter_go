package emby

import "log"

type SystemInfo struct {
	ID                 string `json:"Id"`
	Version            string `json:"Version"`
	WanAddress         string `json:"WanAddress"`
	LocalAddress       string `json:"LocalAddress"`
	HasPendingRestart  bool   `json:"HasPendingRestart"`
	HasUpdateAvailable bool   `json:"HasUpdateAvailable"`
}

func (s *Server) GetServerInfo() *SystemInfo {
	var systemInfo SystemInfo
	err := s.request("GET", "/System/Info", "", &systemInfo)
	if err != nil {
		log.Println("Emby Server - GetServerInfo : " + err.Error())
		return &SystemInfo{
			Version:            "0.0.0",
			HasPendingRestart:  false,
			HasUpdateAvailable: false,
			LocalAddress:       "",
			WanAddress:         "",
		}
	}
	return &systemInfo
}
