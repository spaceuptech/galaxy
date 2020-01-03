package model

// DatabaseEventMessage is the event payload for create, update and delete events
type DatabaseEventMessage struct {
	Data struct {
		DBType string      `json:"db"`
		Col    string      `json:"col"`
		DocID  string      `json:"docId"`
		Doc    interface{} `json:"doc"`
	} `json:"data"`
}
