package entity

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
	Name           string         `json:"Name"`
	LibraryOptions LibraryOptions `json:"LibraryOptions"`
	ItemID         string         `json:"ItemId"`
}

type LibraryMetrics struct {
	Name string
	Size int
}
