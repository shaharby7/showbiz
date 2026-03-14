package provider

import (
"bytes"
"context"
"encoding/json"
"fmt"
"io"
"net/http"
"time"
)

// FakeProvider implements the Provider interface by calling the fakeprovider microservice,
// which manages KubeVirt VirtualMachineInstances.
type FakeProvider struct {
baseURL    string
httpClient *http.Client
}

var _ Provider = (*FakeProvider)(nil)

func NewFakeProvider(baseURL string) *FakeProvider {
return &FakeProvider{
baseURL: baseURL,
httpClient: &http.Client{
Timeout: 10 * time.Second,
},
}
}

func (f *FakeProvider) Name() string {
return "fakeprovider"
}

func (f *FakeProvider) ResourceTypes() []string {
return []string{"machine"}
}

func (f *FakeProvider) ValidateCredentials(_ context.Context, _ map[string]interface{}) error {
// No credentials needed for local kubevirt
return nil
}

// fakeproviderMachine mirrors the fakeprovider service's Machine model.
type fakeproviderMachine struct {
ID        string `json:"id"`
Name      string `json:"name"`
Namespace string `json:"namespace"`
CPU       int    `json:"cpu"`
MemoryMB  int    `json:"memoryMB"`
Image     string `json:"image"`
Status    string `json:"status"`
IP        string `json:"ip,omitempty"`
}

func (f *FakeProvider) CreateResource(ctx context.Context, input *CreateResourceInput) (*ResourceOutput, error) {
cpu := 1
if v, ok := input.Properties["cpu"]; ok {
if n, ok := v.(float64); ok {
cpu = int(n)
}
}
memoryMB := 512
if v, ok := input.Properties["memoryMB"]; ok {
if n, ok := v.(float64); ok {
memoryMB = int(n)
}
}
image := "quay.io/kubevirt/cirros-container-disk-demo"
if v, ok := input.Properties["image"]; ok {
if s, ok := v.(string); ok {
image = s
}
}
namespace := "vmis"
if v, ok := input.Properties["namespace"]; ok {
if s, ok := v.(string); ok {
namespace = s
}
}

body, _ := json.Marshal(map[string]interface{}{
"name":      input.Name,
"namespace": namespace,
"cpu":       cpu,
"memoryMB":  memoryMB,
"image":     image,
})

req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.baseURL+"/v1/machines", bytes.NewReader(body))
if err != nil {
return nil, fmt.Errorf("build request: %w", err)
}
req.Header.Set("Content-Type", "application/json")

resp, err := f.httpClient.Do(req)
if err != nil {
return nil, fmt.Errorf("call fakeprovider: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusCreated {
respBody, _ := io.ReadAll(resp.Body)
return nil, fmt.Errorf("fakeprovider returned %d: %s", resp.StatusCode, string(respBody))
}

var machine fakeproviderMachine
if err := json.NewDecoder(resp.Body).Decode(&machine); err != nil {
return nil, fmt.Errorf("decode response: %w", err)
}

props := map[string]interface{}{
"cpu":                cpu,
"memoryMB":           memoryMB,
"image":              image,
"namespace":          namespace,
"fakeproviderID":     machine.ID,
}

return &ResourceOutput{
ID:         machine.ID,
Type:       "machine",
Name:       input.Name,
Status:     machine.Status,
Properties: props,
}, nil
}

func (f *FakeProvider) GetResource(ctx context.Context, resourceID string) (*ResourceOutput, error) {
req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.baseURL+"/v1/machines/"+resourceID, nil)
if err != nil {
return nil, fmt.Errorf("build request: %w", err)
}

resp, err := f.httpClient.Do(req)
if err != nil {
return nil, fmt.Errorf("call fakeprovider: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode == http.StatusNotFound {
return nil, fmt.Errorf("resource not found")
}
if resp.StatusCode != http.StatusOK {
respBody, _ := io.ReadAll(resp.Body)
return nil, fmt.Errorf("fakeprovider returned %d: %s", resp.StatusCode, string(respBody))
}

var machine fakeproviderMachine
if err := json.NewDecoder(resp.Body).Decode(&machine); err != nil {
return nil, fmt.Errorf("decode response: %w", err)
}

props := map[string]interface{}{
"cpu":            machine.CPU,
"memoryMB":       machine.MemoryMB,
"image":          machine.Image,
"namespace":      machine.Namespace,
"fakeproviderID": machine.ID,
}
if machine.IP != "" {
props["ip"] = machine.IP
}

return &ResourceOutput{
ID:         machine.ID,
Type:       "machine",
Name:       machine.Name,
Status:     machine.Status,
Properties: props,
}, nil
}

func (f *FakeProvider) UpdateResource(ctx context.Context, input *UpdateResourceInput) (*ResourceOutput, error) {
body, _ := json.Marshal(input.Properties)

req, err := http.NewRequestWithContext(ctx, http.MethodPut, f.baseURL+"/v1/machines/"+input.ResourceID, bytes.NewReader(body))
if err != nil {
return nil, fmt.Errorf("build request: %w", err)
}
req.Header.Set("Content-Type", "application/json")

resp, err := f.httpClient.Do(req)
if err != nil {
return nil, fmt.Errorf("call fakeprovider: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode == http.StatusNotFound {
return nil, fmt.Errorf("resource not found")
}

var machine fakeproviderMachine
if err := json.NewDecoder(resp.Body).Decode(&machine); err != nil {
return nil, fmt.Errorf("decode response: %w", err)
}

props := map[string]interface{}{
"cpu":            machine.CPU,
"memoryMB":       machine.MemoryMB,
"image":          machine.Image,
"namespace":      machine.Namespace,
"fakeproviderID": machine.ID,
}

return &ResourceOutput{
ID:         machine.ID,
Type:       "machine",
Name:       machine.Name,
Status:     machine.Status,
Properties: props,
}, nil
}

func (f *FakeProvider) DeleteResource(ctx context.Context, resourceID string) error {
req, err := http.NewRequestWithContext(ctx, http.MethodDelete, f.baseURL+"/v1/machines/"+resourceID, nil)
if err != nil {
return fmt.Errorf("build request: %w", err)
}

resp, err := f.httpClient.Do(req)
if err != nil {
return fmt.Errorf("call fakeprovider: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode == http.StatusNotFound {
return fmt.Errorf("resource not found")
}
if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
return fmt.Errorf("fakeprovider returned %d", resp.StatusCode)
}
return nil
}

func (f *FakeProvider) DetectDrifts(_ context.Context, resources []ResourceExpectedState) ([]DriftReport, error) {
reports := make([]DriftReport, len(resources))
for i, r := range resources {
reports[i] = DriftReport{
ResourceID: r.ResourceID,
Drifted:    false,
}
}
return reports, nil
}
