package model

// DatabaseEventPayload is the payload received from database events
type DatabaseEventPayload struct {
	Data struct {
		Doc *Service `json:"doc"`
	} `json:"data"`
}

// FileStoreEventPayload is the payload received from file store events
type FileStoreEventPayload struct {
	Data *FileStoreData `json:"data"`
}

type FileStoreData struct {
	Path string          `json:"path"`
	Meta *ServiceRequest `json:"meta"`
}

type ServiceRequest struct {
	IsDeploy bool     `json:"isDeploy"`
	Service  *Service `json:"service"`
}
