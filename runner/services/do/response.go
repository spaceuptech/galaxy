package do

import "github.com/digitalocean/godo"

// DOdatabases struct as defined in the response from digitalocean
type DOdatabases struct {
	Databases []godo.Database `json:"databases"`
}

// DOdatabase struct as defined in the response from digitalocean
type DOdatabase struct {
	Database godo.Database `json:"database"`
}

// dbSizeSlug is a map which acts as a lookup table for db sizes
var dbSizeSlug = map[string]struct{}{
	"db-s-1vcpu-1gb":   struct{}{}, // using empty struct consumes no storage space :P
	"db-s-1vcpu-2gb":   struct{}{},
	"db-s-2vcpu-4gb":   struct{}{},
	"db-s-4vcpu-8gb":   struct{}{},
	"db-s-6vcpu-16gb":  struct{}{},
	"db-s-8vcpu-32gb":  struct{}{},
	"db-s-16vcpu-64gb": struct{}{},
}
