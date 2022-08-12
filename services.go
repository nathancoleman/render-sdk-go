package render

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Repo         string    `json:"repo"`
	Name         string    `json:"name"`
	AutoDeploy   bool      `json:"autoDeploy"`
	Branch       string    `json:"branch"`
	CreatedAt    time.Time `json:"createdAt"`
	NotifyOnFail string    `json:"notifyOnFail"`
	OwnerID      string    `json:"ownerId"`
	Slug         string    `json:"slug"`
	Suspended    string    `json:"suspended"`
	Suspenders   []string  `json:"suspenders"`
	UpdatedAt    time.Time `json:"updatedAt"`

	ServiceDetails ServiceDetails `json:"serviceDetails"`
}

type ServiceDetails struct {
	Env                 string    `json:"env"`
	LastSuccessfulRunAt time.Time `json:"lastSuccessfulRunAt"`
	Plan                string    `json:"plan"`
	Region              string    `json:"region"`
	Schedule            string    `json:"schedule,omitempty"`
	HealthCheckPath     string    `json:"healthCheckPath,omitempty"`
	URL                 string    `json:"url,omitempty"`
}

type Services struct {
	client *Client
}

func NewServices(client *Client) *Services {
	return &Services{client}
}

type ListServicesResponseBody []struct {
	Cursor  string  `json:"cursor"`
	Service Service `json:"service"`
}

// List returns all services that are owned by your Render user
// and any teams you're a part of.
//
// https://api-docs.render.com/reference/get-services
func (s *Services) List(ctx context.Context) ([]Service, error) {
	resp, err := s.client.c.ListServicesWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to retrieve services: %d %s", resp.StatusCode(), string(resp.Body))
	}

	services := make([]Service, 0, len(*resp.JSON200))
	for _, service := range *resp.JSON200 {
		services = append(services, Service{
			ID:             service.Service.Id,
			Type:           "",
			Repo:           service.Service.Repo,
			Name:           service.Service.Name,
			AutoDeploy:     service.Service.AutoDeploy == "yes",
			Branch:         service.Service.Branch,
			CreatedAt:      service.Service.CreatedAt,
			NotifyOnFail:   "",
			OwnerID:        "",
			Slug:           "",
			Suspended:      "",
			Suspenders:     nil,
			UpdatedAt:      service.Service.UpdatedAt,
			ServiceDetails: ServiceDetails{},
		})
	}

	return services, nil
}

// Get returns the service with the provided serviceID.
//
// https://api-docs.render.com/reference/get-service
func (s *Services) Get(ctx context.Context, serviceID string) (*Service, error) {
	resp, err := s.client.c.GetServiceWithResponse(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to get service %s: %d %s", serviceID, resp.StatusCode(), string(resp.Body))
	}

	service := *resp.JSON200

	return &Service{
		ID:             service.Id,
		Type:           "",
		Repo:           service.Repo,
		Name:           service.Name,
		AutoDeploy:     service.AutoDeploy == "yes",
		Branch:         service.Branch,
		CreatedAt:      service.CreatedAt,
		NotifyOnFail:   "",
		OwnerID:        "",
		Slug:           "",
		Suspended:      "",
		Suspenders:     nil,
		UpdatedAt:      service.UpdatedAt,
		ServiceDetails: ServiceDetails{},
	}, nil
}
