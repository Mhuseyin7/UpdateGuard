package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/updateguard/updateguard/internal/domain"
)

// HTTPClient talks only to UpdateGuard Agent's allowlisted API; Docker is never exposed to browsers/API users.
type HTTPClient struct { BaseURL, Token string; Client *http.Client }
func (c HTTPClient) call(ctx context.Context, method, path string, in, out any) error {
	var body *bytes.Reader
	if in != nil { raw, err := json.Marshal(in); if err != nil{return err}; body=bytes.NewReader(raw) } else { body=bytes.NewReader(nil) }
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body); if err != nil{return err}
	req.Header.Set("Authorization", "Bearer "+c.Token); req.Header.Set("Content-Type", "application/json")
	r, err := c.Client.Do(req); if err != nil{return err}; defer r.Body.Close()
	if r.StatusCode >= 300 { return fmt.Errorf("agent %s: %s", path, r.Status) }
	return json.NewDecoder(r.Body).Decode(out)
}
func (c HTTPClient) Discover(ctx context.Context) ([]domain.Service,error) { var out []domain.Service; return out,c.call(ctx,"GET","/v1/services",nil,&out) }
func (c HTTPClient) Snapshot(ctx context.Context, id string)(domain.Snapshot,error){var out domain.Snapshot;return out,c.call(ctx,"GET","/v1/services/"+id+"/snapshot",nil,&out)}
func (c HTTPClient) Preflight(ctx context.Context,id string,t domain.Snapshot)error{return c.call(ctx,"POST","/v1/services/"+id+"/preflight",t,&struct{}{})}
func (c HTTPClient) Pull(ctx context.Context,d string)error{return c.call(ctx,"POST","/v1/images/pull",map[string]string{"digest":d},&struct{}{})}
func (c HTTPClient) Recreate(ctx context.Context,id string,s domain.Snapshot)error{return c.call(ctx,"POST","/v1/services/"+id+"/recreate",s,&struct{}{})}
func (c HTTPClient) Verify(ctx context.Context,id string,v domain.Verification)error{return c.call(ctx,"POST","/v1/services/"+id+"/verify",v,&struct{}{})}
func (c HTTPClient) Rollback(ctx context.Context,id string,s domain.Snapshot)error{return c.call(ctx,"POST","/v1/services/"+id+"/rollback",s,&struct{}{})}
func (c HTTPClient) Logs(ctx context.Context,id string,n int)([]string,error){var out []string;return out,c.call(ctx,"GET",fmt.Sprintf("/v1/services/%s/logs?tail=%d",id,n),nil,&out)}
