package entity

type SystemInfo struct {
	ID                 string `json:"Id"`
	Version            string `json:"Version"`
	WanAddress         string `json:"WanAddress"`
	LocalAddress       string `json:"LocalAddress"`
	HasPendingRestart  bool   `json:"HasPendingRestart"`
	HasUpdateAvailable bool   `json:"HasUpdateAvailable"`
}
