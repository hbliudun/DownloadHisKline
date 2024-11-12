package save

// dbbardata
type DBBarData struct {
	Symbol       string  `json:"symbol"`
	Exchange     string  `json:"exchange"`
	Datetime     string  `json:"datetime"`
	Interval     string  `json:"interval"`
	Volume       float64 `json:"volume"`
	Turnover     float64 `json:"turnover"`
	OpenInterest float64 `json:"open_interest"`
	Open         float64 `json:"open"`
	Close        float64 `json:"close"`
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
}

// dbbaroverview
type DBBarOverview struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
	Interval string `json:"interval"`
	Count    int64  `json:"count"`
	Start    string `json:"start"`
	End      string `json:"end"`
}
