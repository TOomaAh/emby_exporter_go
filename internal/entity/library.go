package entity

type Library struct {
	Items            *[]Item `json:"Items"`
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
	LibraryOptions LibraryOptions `json:"LibraryOptions"`
	Name           string         `json:"Name"`
	ItemID         string         `json:"ItemId"`
}
