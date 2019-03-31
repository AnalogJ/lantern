package models

type DbNotify struct {
	Table        string `json:"table"`
	Action	 	 string `json:"action"`
	Data 		map[string]interface{} `json:"data"`
}

