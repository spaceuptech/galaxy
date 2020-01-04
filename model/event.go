package model

// DatabaseEventMessage is the event payload for create, update and delete events
type DatabaseEventMessage struct {
	Data struct {
		Doc *Service `json:"doc"`
	} `json:"data"`
}

type ServiceRequest struct {
	IsDeploy bool     `json:"is_deploy"`
	Service  *Service `json:"service"`
}

// DatabaseEventMessage is the event payload for create, update and delete events
type FileStoreEventMessage struct {
	Data struct {
		Path string         `json:"path"`
		Meta ServiceRequest `json:"meta"`
	} `json:"data"`
}
