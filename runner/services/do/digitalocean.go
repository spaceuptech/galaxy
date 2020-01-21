package do

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"

	"github.com/spaceuptech/galaxy/model"
)

// DigitalOcean is used to manage DO clients
type DigitalOcean struct {
	client *godo.Client
	token  string
	region string
}

// TokenSource -> pat
type TokenSource struct {
	AccessToken string
}

// Token function returns token to create client
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// New creates a digitalOcean client and returns DO object
func New(token string, region string) *DigitalOcean {
	tokenSource := &TokenSource{
		AccessToken: token,
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	return &DigitalOcean{client: client, token: token, region: region}

}

// Apply is used for create/update operations on DB
func (do *DigitalOcean) Apply(ctx context.Context, service *model.ManagedService) error {

	// Check if the request received is a valid request
	if service.ServiceType == "database" {
		if service.DataBase.Replication.ReplicationFactor == service.DataBase.Replication.Instances {

			if service.DataBase.Type != "pg" && service.DataBase.Type != "mysql" && service.DataBase.Type != "redis" {
				return errors.New("invalid Database type provided")
			}

			sizeSlug := "db-s-" + strconv.FormatInt(service.DBResources.CPU, 10) + "vcpu-" + strconv.FormatInt(service.DBResources.Memory, 10) + "gb"
			if _, ok := dbSizeSlug[sizeSlug]; !ok {
				return errors.New("Invalid db size")
			}
			// list all database
			listDB, err := do.listDBsByTag(ctx, getTagName(service.ProjectID, service.ID))
			if err != nil {
				return err
			}
			// Check if database already exists: if it does -> update ; else -> create
			if len(listDB.Databases) == 0 {
				// Create a new DB Cluster
				createRequest := &godo.DatabaseCreateRequest{
					Name:       service.ID,
					EngineSlug: service.DataBase.Type,
					Version:    service.DataBase.DataBaseVersion,
					Region:     do.region,
					SizeSlug:   sizeSlug,
					NumNodes:   service.DataBase.Replication.ReplicationFactor, //Value inclusinve of stand-by nodes
					Tags: []string{
						service.ID,
						service.ProjectID,
					},
				}
				doDB, _, err := do.client.Databases.Create(ctx, createRequest)

				if err != nil {
					return fmt.Errorf("Error creating db cluster: %s", err)
				}

				// Create a new Database-User named 'galaxy' (#password is auto-generated)
				addUserRequest := &godo.DatabaseCreateUserRequest{
					Name: "galaxy",
				}
				_, _, err = do.client.Databases.CreateUser(ctx, doDB.ID, addUserRequest)

				if err != nil {
					return fmt.Errorf("Error creating db user: %s", err)
				}
				return nil
			}

			// Get the database ID
			dbID := listDB.Databases[0].ID
			resizeRequest := &godo.DatabaseResizeRequest{
				SizeSlug: sizeSlug,
				NumNodes: service.DataBase.Replication.ReplicationFactor,
			}
			_, err = do.client.Databases.Resize(ctx, dbID, resizeRequest)
			if err != nil {
				return fmt.Errorf("Error resizing db cluster: %s", err)
			}
			return nil
		}
		return errors.New("Replication Factor MUST be equal to Instances")
	}
	return errors.New("Invalid Service Type received")
}

// Delete is used to delete the database cluster
func (do *DigitalOcean) Delete(ctx context.Context, service *model.ManagedService) error {

	listDB, err := do.listDBsByTag(ctx, getTagName(service.ProjectID, service.ID))
	if err != nil {
		return err
	}
	if len(listDB.Databases) == 0 {
		return fmt.Errorf("database (%s:%s) not found", service.ProjectID, service.ID)
	}
	if _, err := do.client.Databases.Delete(ctx, listDB.Databases[0].ID); err != nil {
		return fmt.Errorf("Error deleting db cluster: %s", err)
	}
	return nil
}

// GetServices returns the user details for the db
func (do *DigitalOcean) GetServices(ctx context.Context, service *model.ManagedService) (*model.GetServiceDetails, error) {

	// Retrieve an existing db cluster by tag..it contains the user as well as the conenction details
	cluster, err := do.listDBsByTag(ctx, getTagName(service.ProjectID, service.ID))
	if err != nil {
		return nil, fmt.Errorf("Error fetching db details: %s", err)
	}
	if len(cluster.Databases) == 0 {
		return nil, fmt.Errorf("database (%s:%s) not found", service.ProjectID, service.ID)
	}

	var uname string
	var password string
	for _, user := range cluster.Databases[0].Users {
		if user.Name == "galaxy" {
			uname = user.Name
			password = user.Password
		}
	}
	// Connection -> Public && PrivateConnection -> Private
	port := cluster.Databases[0].Connection.Port
	pubURI := cluster.Databases[0].Connection.URI
	prvURI := cluster.Databases[0].PrivateConnection.URI
	pubHost := cluster.Databases[0].Connection.Host
	prvHost := cluster.Databases[0].PrivateConnection.Host

	return &model.GetServiceDetails{
		PublicNw: model.PublicNw{
			Username: uname,
			Password: password,
			Port:     port,
			URI:      pubURI,
			Host:     pubHost,
		},
		PrivateNw: model.PrivateNw{
			Username: uname,
			Password: password,
			Port:     port,
			URI:      prvURI,
			Host:     prvHost,
		},
	}, nil
}

// GetAllTech is used to return all possible tech of the vendor selected
func GetAllTech() []string {
	return []string{"mysql", "postgres"}
}
