package do

import (
	"context"
	"fmt"
	"net/http"
)

func getTagName(projectID, serviceID string) string {
	return fmt.Sprintf("%s-%s", projectID, serviceID)
}

func (do *DigitalOcean) listDBsByTag(ctx context.Context, tagName string) (*DOdatabases, error) {
	// list all database
	reqURL := "https://api.digitalocean.com/v2/databases?tag_name=" + tagName //iff not found...create a new dB with tAG as ID/ProjectID
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	// Add the appropriate headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+do.token)

	listDB := new(DOdatabases)
	if _, err := do.client.Do(ctx, req, listDB); err != nil {
		return nil, err
	}
	return listDB, nil
}
