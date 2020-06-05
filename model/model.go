package model

import "time"

// RequestInput struct
type RequestInput struct {
	Date string `json:"date"`
}

// Request struct
type Request struct {
	Date time.Time
}

// ======================== Qraphql Structs ================================ //

// PropaneSales struct
type PropaneSales struct {
	Report struct {
		Date       time.Time      `json:"date"`
		Deliveries map[string]int `json:"deliveries"`
		// Sales      map[string]map[string]float64
		Sales []struct {
			Date  string             `json:"date"`
			Sales map[string]float64 `json:"sales"`
		}
	} `json:"propaneReportDwnld"`
}
