package linear

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/n1x9s/second-brain/backend/internal/domain"
)

const endpoint = "https://api.linear.app/graphql"

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) Client {
	return Client{apiKey: apiKey, client: &http.Client{Timeout: 10 * time.Second}}
}

func (c Client) CreateIssue(ctx context.Context, teamID string, title string, description string) (domain.LinearIssue, error) {
	var out struct {
		IssueCreate struct {
			Success bool        `json:"success"`
			Issue   linearIssue `json:"issue"`
		} `json:"issueCreate"`
	}
	err := c.call(ctx, `mutation IssueCreate($input: IssueCreateInput!) {
		issueCreate(input: $input) { success issue { id identifier title description url state { name } } }
	}`, map[string]any{
		"input": map[string]any{"teamId": teamID, "title": title, "description": description},
	}, &out)
	if err != nil {
		return domain.LinearIssue{}, err
	}
	return out.IssueCreate.Issue.toDomain(), nil
}

func (c Client) UpdateIssue(ctx context.Context, issueID string, title string, description string, stateID string) (domain.LinearIssue, error) {
	input := map[string]any{}
	if title != "" {
		input["title"] = title
	}
	if description != "" {
		input["description"] = description
	}
	if stateID != "" {
		input["stateId"] = stateID
	}
	var out struct {
		IssueUpdate struct {
			Success bool        `json:"success"`
			Issue   linearIssue `json:"issue"`
		} `json:"issueUpdate"`
	}
	err := c.call(ctx, `mutation IssueUpdate($id: String!, $input: IssueUpdateInput!) {
		issueUpdate(id: $id, input: $input) { success issue { id identifier title description url state { name } } }
	}`, map[string]any{"id": issueID, "input": input}, &out)
	if err != nil {
		return domain.LinearIssue{}, err
	}
	return out.IssueUpdate.Issue.toDomain(), nil
}

func (c Client) ListIssues(ctx context.Context, teamID string) ([]domain.LinearIssue, error) {
	var out struct {
		Issues struct {
			Nodes []linearIssue `json:"nodes"`
		} `json:"issues"`
	}
	err := c.call(ctx, `query Issues($teamId: String!) {
		issues(filter: { team: { id: { eq: $teamId } }, completedAt: { null: true } }, first: 50) {
			nodes { id identifier title description url state { name } }
		}
	}`, map[string]any{"teamId": teamID}, &out)
	if err != nil {
		return nil, err
	}
	issues := make([]domain.LinearIssue, 0, len(out.Issues.Nodes))
	for _, issue := range out.Issues.Nodes {
		issues = append(issues, issue.toDomain())
	}
	return issues, nil
}

func (c Client) call(ctx context.Context, query string, variables map[string]any, target any) error {
	if c.apiKey == "" {
		return domain.ErrProviderEmpty
	}
	payload, _ := json.Marshal(map[string]any{"query": query, "variables": variables})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("linear graphql failed: %s", resp.Status)
	}
	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	if len(envelope.Errors) > 0 {
		return fmt.Errorf("linear graphql: %s", envelope.Errors[0].Message)
	}
	return json.Unmarshal(envelope.Data, target)
}

type linearIssue struct {
	ID          string `json:"id"`
	Identifier  string `json:"identifier"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	State       struct {
		Name string `json:"name"`
	} `json:"state"`
}

func (i linearIssue) toDomain() domain.LinearIssue {
	return domain.LinearIssue{ID: i.ID, Identifier: i.Identifier, Title: i.Title, Description: i.Description, State: i.State.Name, URL: i.URL}
}
