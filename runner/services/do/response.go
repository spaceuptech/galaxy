package do

// DOdatabase struct as defined in the response from digitalocean
type DOdatabase struct {
	Databases []DBresponse `json:"databases"`
}

// DBresponse struct defines the responde rxd from the do
type DBresponse struct {
	ID                string
	Name              string
	Engine            string
	Version           string
	Connection        Connection
	PrivateConnection PrivateConnection
	Users             []Users
	DbNames           []string
	NumNodes          int64
	Region            string
	Status            string
	Size              string
	CreatedAt         string
	MaintenanceWindow MaintenanceWindow
	Tags              []string
}

// Connection contains the info needed to access the db cluster
type Connection struct {
	URI      string
	Database string
	Host     string
	Port     int
	User     string
	Password string
	Ssl      bool
}

// PrivateConnection contains the info needed to access the db cluster
type PrivateConnection struct {
	URI      string
	Database string
	Host     string
	Port     int
	User     string
	Password string
	Ssl      bool
}

// Users contains the database users
type Users struct {
	Name     string
	Role     string
	Password string
}

// MaintenanceWindow contains information about any pending maintenance for the db cluster
type MaintenanceWindow struct {
	Day         string
	Hour        string
	Pending     bool
	Description []string
}
