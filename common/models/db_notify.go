package models

type DbNotify struct {
	Table        string `json:"table"`
	Action	 	 string `json:"action"`
	Id 		int `json:"id"`
}

