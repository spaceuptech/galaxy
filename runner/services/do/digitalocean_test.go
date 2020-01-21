package do

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/utils"
)

type httpMock struct {
	activeStep int
	steps      []mockStep
}

type mockStep struct {
	// Fields for request
	method, uri string
	body        interface{}

	// Fields for response
	status  int
	res     interface{}
	wantErr bool
}

func mockDO(steps []mockStep) *godo.Client {
	c := &http.Client{}
	c.Transport = &httpMock{steps: steps}
	return godo.NewClient(c)
}

func (m *httpMock) RoundTrip(r *http.Request) (*http.Response, error) {
	// Get the active step
	if m.activeStep == len(m.steps) {
		return nil, errors.New("not enough steps provided")
	}
	step := m.steps[m.activeStep]
	m.activeStep++

	if r.Body != nil && r.Body != http.NoBody && step.body == nil {
		return nil, errors.New("request sent a body which is not required")
	}

	if step.body != nil {
		if r.Body == nil {
			return nil, fmt.Errorf("invalid request body")
		}
		var reqBody interface{}
		_ = json.NewDecoder(r.Body).Decode(&reqBody)
		if reflect.DeepEqual(reqBody, step.body) {
			return nil, fmt.Errorf("invalid body sent - %v", reqBody)
		}
		defer utils.CloseReaderCloser(r.Body)
	}

	if step.method != r.Method {
		return nil, fmt.Errorf("invalid request method - got %s; wanted %s", r.Method, step.method)
	}

	if step.uri != r.URL.RequestURI() {
		return nil, fmt.Errorf("invalid request uri - got %s; wanted %s", r.URL.RequestURI(), step.uri)
	}

	if step.wantErr {
		return nil, errors.New("you wanted it, so you got it")
	}

	reader, writer := io.Pipe()
	go func() {
		json.NewEncoder(writer).Encode(step.res)
		writer.Close()
	}()
	return &http.Response{StatusCode: step.status, Body: reader}, nil
}

func TestDigitalOcean_Delete(t *testing.T) {

	tests := []struct {
		name    string
		steps   []mockStep
		service *model.ManagedService
		wantErr bool
	}{
		{
			name: "check if request sent is right", wantErr: false,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p1-s1", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{{ID: "10"}}}},
				{method: http.MethodDelete, uri: "/v2/databases/10", body: nil, status: 204, res: nil},
			},
			service: &model.ManagedService{ID: "s1", ProjectID: "p1"},
		},
		{
			name: "empty databases reponse", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p1-s1", body: nil, status: 400, res: DOdatabases{Databases: []godo.Database{}}}, // should fail here
				{method: http.MethodDelete, uri: "/v2/databases/10", body: nil, status: 204, res: nil},
			},
			service: &model.ManagedService{ID: "s1", ProjectID: "p1"},
		}, {
			name: "invalid list response", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p1-s1", body: nil, status: 400, wantErr: true}, // should fail here
				{method: http.MethodDelete, uri: "/v2/databases/10", body: nil, status: 204, res: nil},
			},
			service: &model.ManagedService{ID: "s1", ProjectID: "p1"},
		},
		{
			name: "database not found", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p1-s1", body: nil, status: 204, res: DOdatabases{Databases: []godo.Database{}}},
				{method: http.MethodDelete, uri: "/v2/databases/10", body: nil, status: 400, res: nil, wantErr: true},
			},
			service: &model.ManagedService{ID: "s1", ProjectID: "p1"},
		},
		{
			name: "error deleting database", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p1-s1", body: nil, status: 204, res: DOdatabases{Databases: []godo.Database{{ID: "10"}}}},
				{method: http.MethodDelete, uri: "/v2/databases/11", body: nil, status: 400, res: nil, wantErr: true}, // should fail here
			},
			service: &model.ManagedService{ID: "s1", ProjectID: "p1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			do := &DigitalOcean{client: mockDO(tt.steps), region: "nyc"}
			if err := do.Delete(context.Background(), tt.service); (err != nil) != tt.wantErr {
				t.Errorf("DigitalOcean.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDigitalOcean_Apply(t *testing.T) {
	type args struct {
		ctx     context.Context
		service *model.ManagedService
	}
	tests := []struct {
		name    string
		steps   []mockStep
		service *model.ManagedService
		wantErr bool
	}{
		// test: Create successfully
		{
			name: "check creating db", wantErr: false,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s2"}}},
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 200, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: Create -> empty response check
		{
			name: "db already existing error check", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{{ID: "s2"}}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: http.StatusBadRequest, res: DOdatabase{godo.Database{ID: "s2"}}}, // should fail here
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 200, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: Create -> invalid/empty response while listing db
		{
			name: "invalid/empty listDB resp", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 400, wantErr: true}, // should fail here
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s2"}}},
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 200, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: Create -> empty response after creating a db cluster check
		{
			name: "empty response after creating db cluster", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{{ID: "s2"}}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 400, res: nil}, // should fail here
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 200, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: Create -> empty response after creating a db user check
		{
			name: "empty response after creating a user cluster", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{{ID: "s2"}}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s2"}}},
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 400, res: nil}, // should fail here
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: invalid slug-size
		{
			name: "invalid-slug-size", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s2"}}}, // should fail here
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 200, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      4,
					Memory:   4,
					IsShared: false,
				},
			},
		},
		// test: error creating db cluster due to invalid create request
		{
			name: "invalid create db Request", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s2"}}, wantErr: true}, // should fail here
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 200, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pgggg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: error creating db user due to invalid create request
		{
			name: "invalid create user Request", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s"}}},
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 400, res: DOdatabase{godo.Database{}}, wantErr: true}, // should fail here
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: repln factor != instances
		{
			name: "replication factor != instances", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s2"}}, wantErr: true}, // should fail here
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 400, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 4,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: invalid service type
		{
			name: "invalid serviceType", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{}}},
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s2"}}, wantErr: true}, // should fail here
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 400, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "db",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 4,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// test: error fetching db details due to invalid id or project id
		{
			name: "error fetching db", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p-s", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{}}, wantErr: true}, // should fail here
				{method: http.MethodPost, uri: "/v2/databases", body: godo.DatabaseCreateRequest{}, status: 200, res: DOdatabase{godo.Database{ID: "s2"}}},
				{method: http.MethodPost, uri: "/v2/databases/s2/users", body: godo.DatabaseCreateUserRequest{}, status: 400, res: DOdatabase{godo.Database{}}},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 4,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},

		//Update case
		{
			name: "check updating db", wantErr: false,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{{ID: "10"}}}},
				{method: http.MethodPut, uri: "/v2/databases/10/resize", body: godo.DatabaseResizeRequest{}, status: 200, res: nil},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		//update -> empty db check
		{
			name: "db does not exist check", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 400, res: DOdatabases{Databases: []godo.Database{}}}, // should fail here
				{method: http.MethodPut, uri: "/v2/databases/10/resize", body: godo.DatabaseResizeRequest{}, status: 200, res: nil},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
		// empty list response
		{
			name: "db does not exist check", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 400, wantErr: true}, // should fail here
				{method: http.MethodPut, uri: "/v2/databases/10/resize", body: godo.DatabaseResizeRequest{}, status: 200, res: nil},
			},
			service: &model.ManagedService{
				ServiceType: "database",
				ID:          "s2",
				ProjectID:   "p2",
				Name:        "creatingGalaxy",
				DataBase: &model.DataBase{
					Type:            "pg",
					DataBaseVersion: "1",
					Replication: model.Replication{
						ReplicationFactor: 2,
						Instances:         2,
					},
				},
				DBResources: model.DBResources{
					Disk:     10,
					CPU:      1,
					Memory:   1,
					IsShared: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			do := &DigitalOcean{client: mockDO(tt.steps), region: "nyc"}
			if err := do.Apply(context.Background(), tt.service); (err != nil) != tt.wantErr {
				t.Errorf("DigitalOcean.Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDigitalOcean_GetServices(t *testing.T) {
	type args struct {
		ctx     context.Context
		service *model.ManagedService
	}
	tests := []struct {
		name    string
		steps   []mockStep
		service *model.ManagedService
		wantErr bool
	}{
		// Right request sent
		{
			name: "check if right request is sent", wantErr: false,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: DOdatabases{Databases: []godo.Database{{
					ID: "10",
					Users: []godo.DatabaseUser{{
						Name:     "galaxy",
						Password: "spaceUp",
					}},
					Connection: &godo.DatabaseConnection{
						URI:  "localhost",
						Port: 1000,
						Host: "space",
					},
					PrivateConnection: &godo.DatabaseConnection{
						URI:  "localhost",
						Port: 1000,
						Host: "space",
					},
				}}}},
			},
			service: &model.ManagedService{ID: "s2", ProjectID: "p2"},
		},

		// Empty db retrieved
		{
			name: "empty db retrieved", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 200, res: nil},
			},
			service: &model.ManagedService{ID: "s2", ProjectID: "p2"},
		},
		// test: error fetching db details
		{
			name: "error fetching db details", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2-s2", body: nil, status: 400, res: DOdatabases{Databases: []godo.Database{{
					Users: []godo.DatabaseUser{{
						Name:     "galaxy",
						Password: "spaceUp",
					}},
					Connection: &godo.DatabaseConnection{
						URI:  "localhost",
						Port: 1000,
						Host: "space",
					},
					PrivateConnection: &godo.DatabaseConnection{
						URI:  "localhost",
						Port: 1000,
						Host: "space",
					},
				}}}},
			},
			service: &model.ManagedService{ID: "s2", ProjectID: "p2"},
		},
		// test: error sending httpRequest to list database cluster
		{
			name: "error sending httpRequest to list database cluster", wantErr: true,
			steps: []mockStep{
				{method: http.MethodGet, uri: "/v2/databases?tag_name=p2", body: nil, status: 400, res: DOdatabases{Databases: []godo.Database{{
					ID: "10",
					Users: []godo.DatabaseUser{{
						Name:     "galaxy",
						Password: "spaceUp",
					}},
					Connection: &godo.DatabaseConnection{
						URI:  "localhost",
						Port: 1000,
						Host: "space",
					},
					PrivateConnection: &godo.DatabaseConnection{
						URI:  "localhost",
						Port: 1000,
						Host: "space",
					},
				}}}},
			},
			service: &model.ManagedService{ID: "s2", ProjectID: "p2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			do := &DigitalOcean{client: mockDO(tt.steps), region: "nyc"}
			if _, err := do.GetServices(context.Background(), tt.service); (err != nil) != tt.wantErr {
				t.Errorf("DigitalOcean.GetServices() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
