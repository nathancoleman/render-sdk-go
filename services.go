package render

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Service struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Repo         string    `json:"repo"`
	Name         string    `json:"name"`
	AutoDeploy   string    `json:"autoDeploy"`
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
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/services", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := ListServicesResponseBody{}
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	services := make([]Service, 0, len(respBody))
	for _, item := range respBody {
		services = append(services, item.Service)
	}

	return services, nil
}

// Retrieve returns the service with the provided serviceID.
//
// https://api-docs.render.com/reference/get-service
func (s *Services) Retrieve(ctx context.Context, serviceID string) (*Service, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/services/"+serviceID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	service := &Service{}
	if err = json.NewDecoder(resp.Body).Decode(service); err != nil {
		return nil, err
	}

	return service, nil
}
