package emby

import "log"

type Library struct {
	Items            []*Item `json:"Items"`
	TotalRecordCount int     `json:"TotalRecordCount"`
}

type Item struct {
	Name string `json:"Name"`
}

type LibraryInfo struct {
	LibraryItem      []LibraryItem `json:"Items"`
	TotalRecordCount int           `json:"TotalRecordCount"`
}

type LibraryOptions struct {
	ContentType string `json:"ContentType"`
}

type LibraryItem struct {
	Name               string         `json:"Name"`
	Locations          []interface{}  `json:"Locations"`
	CollectionType     string         `json:"CollectionType"`
	LibraryOptions     LibraryOptions `json:"LibraryOptions"`
	ItemID             string         `json:"ItemId"`
	PrimaryImageItemID string         `json:"PrimaryImageItemId"`
	RefreshStatus      string         `json:"RefreshStatus"`
}

type LibraryMetrics struct {
	Name string
	Size int
}

func (s *Server) GetLibrary() *LibraryInfo {
	var library LibraryInfo
	err := s.request("GET", "/Library/VirtualFolders/Query", "", &library)

	if err != nil {
		log.Println("Emby Server - GetLibrary : " + err.Error())
		return &library
	}

	return &library
}

func (s *Server) GetLibrarySize(libraryItem *LibraryItem) int {
	var librarySize int
	var library Library
	err := s.request("GET",
		//Ok I need minimum information. Only one Item and api returns the total number of items
		"/Users/"+
			s.UserID+
			"/Items?IncludeItemTypes=Movie&Recursive=true&Fields=BasicSyncInfo&EnableImageTypes=Primary&ParentId="+
			libraryItem.ItemID+"&Limit=1&IncludeItemTypes="+includeType[libraryItem.LibraryOptions.ContentType], "", &library)

	if err != nil {
		log.Println("Cannot get library size, maybe your server is unreachable or your user is not allowed to access this library")
		return 0
	}

	librarySize = library.TotalRecordCount

	return librarySize
}
