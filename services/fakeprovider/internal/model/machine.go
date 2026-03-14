package model

import "time"

type Machine struct {
ID        string    `json:"id"`
Name      string    `json:"name"`
Namespace string    `json:"namespace"`
CPU       int       `json:"cpu"`
MemoryMB  int       `json:"memoryMB"`
Image     string    `json:"image"`
Status    string    `json:"status"`
IP        string    `json:"ip,omitempty"`
CreatedAt time.Time `json:"createdAt"`
UpdatedAt time.Time `json:"updatedAt"`
}

const (
StatusInitialized  = "Initialized"
StatusProvisioning = "Provisioning"
StatusReady        = "Ready"
StatusFailed       = "Failed"
StatusDeleted      = "Deleted"
)
