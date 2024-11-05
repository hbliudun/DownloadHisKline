package save

// dbbaroverview
type DBBarOverview struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
	Interval string `json:"interval"`
	Count    int64  `json:"count"`
	Start    string `json:"start"`
	End      string `json:"end"`
}
