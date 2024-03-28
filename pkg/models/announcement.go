package models

type Announcement struct {
	Header string `json:"header"`
	Info   string `json:"info"`
	Photo  string `json:"photo_href"`
	Cost   uint64 `json:"cost"`
}

type Announcements struct {
	Count         uint64
	Announcements []Announcement
}
