package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is a thin HTTP client for FireWeave's bearer-auth /v1 management plane.
type Client struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
}

func New(endpoint, apiKey string) *Client {
	endpoint = strings.TrimRight(endpoint, "/")
	return &Client{
		endpoint: endpoint,
		apiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("fireweave api error %d: %s", e.StatusCode, e.Body)
}

func (c *Client) do(ctx context.Context, method, path string, in any, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: string(raw)}
	}

	if out == nil || len(raw) == 0 || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	return json.Unmarshal(raw, out)
}

// --- Projects ----------------------------------------------------------------

type Project struct {
	ProjectID   string  `json:"projectId"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type projectEnvelope struct {
	Project Project `json:"project"`
}

type projectsEnvelope struct {
	Projects []Project `json:"projects"`
}

type CreateProjectInput struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
}

type UpdateProjectInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	var out projectsEnvelope
	if err := c.do(ctx, http.MethodGet, "/v1/projects", nil, &out); err != nil {
		return nil, err
	}
	return out.Projects, nil
}

func (c *Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	var out projectEnvelope
	if err := c.do(ctx, http.MethodGet, "/v1/projects/"+projectID, nil, &out); err != nil {
		return nil, err
	}
	return &out.Project, nil
}

func (c *Client) CreateProject(ctx context.Context, in CreateProjectInput) (*Project, error) {
	var out projectEnvelope
	if err := c.do(ctx, http.MethodPost, "/v1/projects", in, &out); err != nil {
		return nil, err
	}
	return &out.Project, nil
}

func (c *Client) UpdateProject(ctx context.Context, projectID string, in UpdateProjectInput) (*Project, error) {
	var out projectEnvelope
	if err := c.do(ctx, http.MethodPatch, "/v1/projects/"+projectID, in, &out); err != nil {
		return nil, err
	}
	return &out.Project, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/projects/"+projectID, nil, nil)
}

// --- Environments ------------------------------------------------------------

type RefRule struct {
	Kind  string  `json:"kind"`
	Value string  `json:"value"`
	Repo  *string `json:"repo,omitempty"`
}

type Environment struct {
	EnvID           string    `json:"envId"`
	ProjectID       string    `json:"projectId"`
	EnvironmentName string    `json:"environmentName"`
	Slug            string    `json:"slug"`
	DisplayName     string    `json:"displayName"`
	PromotionRank   int       `json:"promotionRank"`
	IsDefault       bool      `json:"isDefault"`
	BranchRules     []RefRule `json:"branchRules"`
	TagRules        []RefRule `json:"tagRules"`
	WebhookAliases  []string  `json:"webhookAliases"`
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string    `json:"updatedAt"`
}

type environmentEnvelope struct {
	Environment Environment `json:"environment"`
}

type listEnvironmentsEnvelope struct {
	Environments []Environment `json:"environments"`
}

type CreateEnvironmentInput struct {
	Slug           string    `json:"slug"`
	DisplayName    string    `json:"displayName"`
	BranchRules    []RefRule `json:"branchRules,omitempty"`
	TagRules       []RefRule `json:"tagRules,omitempty"`
	WebhookAliases []string  `json:"webhookAliases,omitempty"`
	IsDefault      *bool     `json:"isDefault,omitempty"`
}

type UpdateEnvironmentInput struct {
	DisplayName    *string   `json:"displayName,omitempty"`
	BranchRules    []RefRule `json:"branchRules,omitempty"`
	TagRules       []RefRule `json:"tagRules,omitempty"`
	WebhookAliases []string  `json:"webhookAliases,omitempty"`
	IsDefault      *bool     `json:"isDefault,omitempty"`
}

func (c *Client) ListEnvironments(ctx context.Context, projectID string) ([]Environment, error) {
	// List endpoint returns a richer CLI shape; map down to Environment when possible.
	var raw map[string]json.RawMessage
	if err := c.do(ctx, http.MethodGet, "/v1/projects/"+projectID+"/environments", nil, &raw); err != nil {
		return nil, err
	}
	var envs []Environment
	if payload, ok := raw["environments"]; ok {
		if err := json.Unmarshal(payload, &envs); err != nil {
			// Fall back to list-shape fields (slug/displayName only).
			var listShape []struct {
				Slug        string    `json:"slug"`
				DisplayName string    `json:"displayName"`
				IsDefault   bool      `json:"isDefault"`
				BranchRules []RefRule `json:"branchRules"`
				TagRules    []RefRule `json:"tagRules"`
			}
			if err2 := json.Unmarshal(payload, &listShape); err2 != nil {
				return nil, err
			}
			for _, e := range listShape {
				envs = append(envs, Environment{
					ProjectID:   projectID,
					Slug:        e.Slug,
					DisplayName: e.DisplayName,
					IsDefault:   e.IsDefault,
					BranchRules: e.BranchRules,
					TagRules:    e.TagRules,
				})
			}
		}
	}
	return envs, nil
}

func (c *Client) GetEnvironment(ctx context.Context, projectID, envID string) (*Environment, error) {
	var out environmentEnvelope
	path := fmt.Sprintf("/v1/projects/%s/environments/%s", projectID, envID)
	if err := c.do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out.Environment, nil
}

func (c *Client) CreateEnvironment(ctx context.Context, projectID string, in CreateEnvironmentInput) (*Environment, error) {
	var out environmentEnvelope
	path := fmt.Sprintf("/v1/projects/%s/environments", projectID)
	if err := c.do(ctx, http.MethodPost, path, in, &out); err != nil {
		return nil, err
	}
	return &out.Environment, nil
}

func (c *Client) UpdateEnvironment(ctx context.Context, projectID, envID string, in UpdateEnvironmentInput) (*Environment, error) {
	var out environmentEnvelope
	path := fmt.Sprintf("/v1/projects/%s/environments/%s", projectID, envID)
	if err := c.do(ctx, http.MethodPatch, path, in, &out); err != nil {
		return nil, err
	}
	return &out.Environment, nil
}

func (c *Client) DeleteEnvironment(ctx context.Context, projectID, envID string) error {
	path := fmt.Sprintf("/v1/projects/%s/environments/%s", projectID, envID)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// --- Org members -------------------------------------------------------------

// OrgMember is a membership row returned by the organizations members API.
type OrgMember struct {
	ID        string  `json:"id"`
	UserID    string  `json:"userId"`
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	Role      string  `json:"role"`
	CreatedAt string  `json:"createdAt"`
}

type orgAPISuccess[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type SetOrgMemberRoleInput struct {
	Role string `json:"role"`
}

type SetOrgMemberRoleResult struct {
	MemberID string `json:"memberId"`
	UserID   string `json:"userId"`
	Role     string `json:"role"`
}

func (c *Client) ListOrgMembers(ctx context.Context, orgID string) ([]OrgMember, error) {
	var out orgAPISuccess[[]OrgMember]
	path := fmt.Sprintf("/api/organizations/%s/members", orgID)
	if err := c.do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

func (c *Client) GetOrgMember(ctx context.Context, orgID, userID string) (*OrgMember, error) {
	members, err := c.ListOrgMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}
	for i := range members {
		if members[i].UserID == userID {
			return &members[i], nil
		}
	}
	return nil, &APIError{StatusCode: http.StatusNotFound, Body: "member not found in organization"}
}

func (c *Client) SetOrgMemberRole(ctx context.Context, orgID, userID string, in SetOrgMemberRoleInput) (*SetOrgMemberRoleResult, error) {
	var out orgAPISuccess[SetOrgMemberRoleResult]
	path := fmt.Sprintf("/api/organizations/%s/members/%s/role", orgID, userID)
	if err := c.do(ctx, http.MethodPut, path, in, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) RemoveOrgMember(ctx context.Context, orgID, userID string) error {
	path := fmt.Sprintf("/api/organizations/%s/members/%s", orgID, userID)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// IsNotFound reports whether err is an HTTP 404 from the FireWeave API.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}
