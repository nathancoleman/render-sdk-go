package render

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Service struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Repo       string    `json:"repo"`
	Name       string    `json:"name"`
	AutoDeploy string    `json:"autoDeploy"`
	Branch     string    `json:"branch"`
	CreatedAt  time.Time `json:"createdAt"`
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
