package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FireWeave-HQ/terraform-provider-fireweave/internal/client"
)

func TestCreateAndGetProject(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fw_org_test" {
			t.Fatalf("missing auth header")
		}
		switch r.Method {
		case http.MethodPost:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"project": map[string]any{
					"projectId": "p1",
					"name":      "Demo",
					"slug":      "demo",
					"status":    "active",
					"createdAt": "2026-01-01T00:00:00Z",
					"updatedAt": "2026-01-01T00:00:00Z",
				},
			})
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"projects": []map[string]any{
					{
						"projectId": "p1",
						"name":      "Demo",
						"slug":      "demo",
						"createdAt": "2026-01-01T00:00:00Z",
						"updatedAt": "2026-01-01T00:00:00Z",
					},
				},
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/projects/p1", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"project": map[string]any{
				"projectId": "p1",
				"name":      "Demo",
				"slug":      "demo",
				"status":    "active",
				"createdAt": "2026-01-01T00:00:00Z",
				"updatedAt": "2026-01-01T00:00:00Z",
			},
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := client.New(srv.URL, "fw_org_test")
	created, err := c.CreateProject(context.Background(), client.CreateProjectInput{
		Name: "Demo",
		Slug: "demo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ProjectID != "p1" {
		t.Fatalf("got %s", created.ProjectID)
	}

	got, err := c.GetProject(context.Background(), "p1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Slug != "demo" {
		t.Fatalf("got slug %s", got.Slug)
	}
}
