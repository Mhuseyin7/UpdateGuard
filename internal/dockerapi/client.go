// Package dockerapi implements the small, allowlisted Docker Engine API surface
// used by UpdateGuard Agent. It is intentionally not a general Docker proxy.
package dockerapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/updateguard/updateguard/internal/domain"
)

const defaultSocket = "/var/run/docker.sock"

type Client struct {
	http *http.Client
}

func New(socket string) *Client {
	if socket == "" {
		socket = defaultSocket
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, "unix", socket)
		},
	}
	return &Client{http: &http.Client{Transport: transport, Timeout: 30 * time.Second}}
}

func (c *Client) request(ctx context.Context, method, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, "http://docker"+path, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("docker engine request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("docker engine %s: %s", path, resp.Status)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) Ping(ctx context.Context) error {
	return c.request(ctx, http.MethodGet, "/_ping", nil)
}

type containerSummary struct {
	ID      string            `json:"Id"`
	Image   string            `json:"Image"`
	ImageID string            `json:"ImageID"`
	Names   []string          `json:"Names"`
	State   string            `json:"State"`
	Labels  map[string]string `json:"Labels"`
}

// Discover returns one service per managed container. Compose labels provide an
// unambiguous project/service relationship without inferring it from names.
func (c *Client) Discover(ctx context.Context, hostID string) ([]domain.Service, error) {
	var containers []containerSummary
	if err := c.request(ctx, http.MethodGet, "/containers/json?all=1", &containers); err != nil {
		return nil, err
	}
	services := make([]domain.Service, 0, len(containers))
	for _, container := range containers {
		name := strings.TrimPrefix(first(container.Names), "/")
		stack := container.Labels["com.docker.compose.project"]
		service := container.Labels["com.docker.compose.service"]
		if service == "" {
			service = name
		}
		services = append(services, domain.Service{
			ID:            container.ID,
			HostID:        hostID,
			Stack:         stack,
			Name:          service,
			Image:         container.Image,
			CurrentDigest: container.ImageID,
			Health:        healthFromState(container.State),
			Policy:        domain.Manual,
		})
	}
	return services, nil
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func healthFromState(state string) string {
	if state == "running" {
		return "RUNNING"
	}
	return strings.ToUpper(state)
}

func ImageInspectPath(reference string) string {
	return "/images/" + url.PathEscape(reference) + "/json"
}
