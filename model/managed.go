package model

// ManagedService describes the database requirements
type ManagedService struct {
	Name        string      `json:"name" yaml:"name"`
	ID          string      `json:"Id" yaml:"Id"`
	ProjectID   string      `json:"projectId" yaml:"projectId"`
	DBResources DBResources `json:"resources" yaml:"resources"`
	Shards      int         `json:"shards" yaml:"shards"`
	ServiceType string      `json:"serviceType" yaml:"serviceType"`
	DataBase    *DataBase   `json:"dataBase" yaml:"dataBase"`
}

// DataBase describes the type of database
type DataBase struct {
	Type            string      `json:"type" yaml:"type"` //The possible values are: "pg" for PostgreSQL, "mysql" for MySQL, and "redis" for Redis.
	DataBaseVersion string      `json:"dbVersion" yaml:"dbVersion"`
	Replication     Replication `json:"replication" yaml:"replication"`
}

// Replication describes the number of standby instances required.
type Replication struct {
	ReplicationFactor int `json:"replicationFactor" yaml:"replicationFactor"` // primary+ standby instances
	Instances         int `json:"instances" yaml:"instances"`
}

// DBResources describes the resources required
type DBResources struct {
	Disk     int64 `json:"disk" yaml:"disk"`
	CPU      int64 `json:"cpu" yaml:"cpu"`
	Memory   int64 `json:"memory" yaml:"memory"`
	IsShared bool  `json:"isShared" yaml:"isShared"`
}

// GetServiceDetails contains the private and public network connection details
type GetServiceDetails struct {
	PublicNw  PublicNw
	PrivateNw PrivateNw
}

// PublicNw contains the nfo related to public connection
type PublicNw struct {
	Username string
	Password string
	Host     string
	Port     int
	URI      string
}

// PrivateNw contains the nfo related to private connection
type PrivateNw struct {
	Username string
	Password string
	Host     string
	Port     int
	URI      string
}
