package main

type Config struct {
	Content []Content `json:"content"`
	Result  Result    `json:"result"`
}
type ImgLayer struct {
	BlockX    int    `json:"BlockX"`
	BlockY    int    `json:"BlockY"`
	ImgFormat string `json:"ImgFormat"`
	ImgLevel  int    `json:"ImgLevel"`
}
type Panos struct {
	Dir   int    `json:"DIR"`
	Order int    `json:"Order"`
	Pid   string `json:"PID"`
	Type  string `json:"Type"`
	X     int    `json:"X"`
	Y     int    `json:"Y"`
}
type Roads struct {
	ID        string  `json:"ID"`
	IsCurrent int     `json:"IsCurrent"`
	Name      string  `json:"Name"`
	Panos     []Panos `json:"Panos"`
	Width     int     `json:"Width"`
}
type TimeLine struct {
	ID        string `json:"ID"`
	IsCurrent int    `json:"IsCurrent"`
	Time      string `json:"Time"`
	TimeDir   string `json:"TimeDir"`
	TimeLine  string `json:"TimeLine"`
	Year      string `json:"Year"`
}
type Content struct {
	Admission    string     `json:"Admission"`
	Date         string     `json:"Date"`
	DeviceHeight float64    `json:"DeviceHeight"`
	Enters       []any      `json:"Enters"`
	FileTag      string     `json:"FileTag"`
	HasDepth     int        `json:"HasDepth"`
	HasTopo      int        `json:"HasTopo"`
	Heading      float64    `json:"Heading"`
	ID           string     `json:"ID"`
	ImgLayer     []ImgLayer `json:"ImgLayer"`
	Inters       []any      `json:"Inters"`
	LayerCount   int        `json:"LayerCount"`
	Links        []any      `json:"Links"`
	Mode         string     `json:"Mode"`
	MoveDir      float64    `json:"MoveDir"`
	NorthDir     float64    `json:"NorthDir"`
	Obsolete     int        `json:"Obsolete"`
	Photos       []any      `json:"Photos"`
	Pitch        int        `json:"Pitch"`
	Provider     int        `json:"Provider"`
	Rx           int        `json:"RX"`
	Ry           int        `json:"RY"`
	Rname        string     `json:"Rname"`
	Roads        []Roads    `json:"Roads"`
	Roll         int        `json:"Roll"`
	Source       string     `json:"Source"`
	SwitchID     []any      `json:"SwitchID"`
	Time         string     `json:"Time"`
	TimeLine     []TimeLine `json:"TimeLine"`
	Type         string     `json:"Type"`
	UserID       string     `json:"UserID"`
	Username     string     `json:"Username"`
	Version      string     `json:"Version"`
	X            int        `json:"X"`
	Y            int        `json:"Y"`
	Z            float64    `json:"Z"`
	FormatV      string     `json:"format_v"`
	Plane        string     `json:"plane"`
	Procdate     string     `json:"procdate"`
}
type Result struct {
	Error int `json:"error"`
}
