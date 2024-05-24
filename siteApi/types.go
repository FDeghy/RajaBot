package siteapi

type Routes []*Route

type Route struct {
	ID   string
	Name string
	Src  string
	Dst  string
}

type Train struct {
	ID                  int    `json:"iD"`
	TrainID             int    `json:"trainID"`
	SourceStationName   string `json:"sourceStationName"`
	TargetStationName   string `json:"targetStationName"`
	Date                string `json:"date"`
	StartTime           string `json:"startTime"`
	EndTime             string `json:"endTime"`
	Number              int    `json:"number"`
	SeatRest            int    `json:"seatRest"`
	SeatPrice           int    `json:"seatPrice"`
	WagonHelp           string `json:"wagonHelp"`
	CompartmentHelp     string `json:"compartmentHelp"`
	CompanyID           int    `json:"companyID"`
	CompanyName         string `json:"companyName"`
	Util                string `json:"util"`
	Path                string `json:"path"`
	SourceStationID     int    `json:"sourceStationID"`
	TargetStationID     int    `json:"targetStationID"`
	DegreeID            int    `json:"degreeID"`
	IsCompartment       int    `json:"isCompartment"`
	IsSallon            int    `json:"isSallon"`
	CompartmentCapacity int    `json:"compartmentCapacity"`
	IsGo                int    `json:"isGo"`
}
