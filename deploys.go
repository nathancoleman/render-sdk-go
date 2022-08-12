package render

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type Deploy struct {
	ID         string     `json:"id"`
	Commit     *Commit    `json:"commit"`
	Status     *string    `json:"status"`
	CreatedAt  *time.Time `json:"createdAt"`
	FinishedAt *time.Time `json:"finishedAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

type Deploys struct {
	client    *Client
	serviceID string
}

func NewDeploys(client *Client, serviceID string) *Deploys {
	return &Deploys{client, serviceID}
}

func (d *Deploys) List(ctx context.Context) ([]Deploy, error) {
	resp, err := d.client.c.ListDeploysWithResponse(ctx, d.serviceID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to retrieve deploys: %d %s", resp.StatusCode(), string(resp.Body))
	}

	deploys := make([]Deploy, 0, len(*resp.JSON200))
	for _, deploy := range *resp.JSON200 {
		var commit *Commit
		if deploy.Deploy.Commit != nil {
			commit = &Commit{
				ID:        (*deploy.Deploy.Commit).Id,
				Message:   (*deploy.Deploy.Commit).Message,
				CreatedAt: (*deploy.Deploy.Commit).CreatedAt,
			}
		}

		var status string
		if deploy.Deploy.Status != nil {
			status = string(*deploy.Deploy.Status)
		}

		deploys = append(deploys, Deploy{
			ID:         deploy.Deploy.Id,
			Commit:     commit,
			Status:     &status,
			CreatedAt:  deploy.Deploy.CreatedAt,
			FinishedAt: deploy.Deploy.FinishedAt,
			UpdatedAt:  deploy.Deploy.UpdatedAt,
		})
	}

	return deploys, nil
}

func (d *Deploys) Get(ctx context.Context, deployID string) (*Deploy, error) {
	resp, err := d.client.c.GetDeployWithResponse(ctx, d.serviceID, deployID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to get deploy %s: %d %s", deployID, resp.StatusCode(), string(resp.Body))
	}

	deploy := *resp.JSON200

	var commit *Commit
	if deploy.Commit != nil {
		commit = &Commit{
			ID:        (*deploy.Commit).Id,
			Message:   (*deploy.Commit).Message,
			CreatedAt: (*deploy.Commit).CreatedAt,
		}
	}

	var status string
	if deploy.Status != nil {
		status = string(*deploy.Status)
	}

	return &Deploy{
		ID:         deploy.Id,
		Commit:     commit,
		Status:     &status,
		CreatedAt:  deploy.CreatedAt,
		FinishedAt: deploy.FinishedAt,
		UpdatedAt:  deploy.UpdatedAt,
	}, nil
}
