package kubevirt

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var vmiGVR = schema.GroupVersionResource{
	Group:    "kubevirt.io",
	Version:  "v1",
	Resource: "virtualmachineinstances",
}

type Client struct {
	dyn dynamic.Interface
}

func NewClient(kubeconfig string) (*Client, error) {
	var cfg *rest.Config
	var err error

	if kubeconfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		cfg, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("build kube config: %w", err)
	}

	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create dynamic client: %w", err)
	}

	return &Client{dyn: dyn}, nil
}

func (c *Client) CreateVMI(ctx context.Context, name, namespace string, cpu, memoryMB int, image string) error {
	vmi := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kubevirt.io/v1",
			"kind":       "VirtualMachineInstance",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"domain": map[string]interface{}{
					"resources": map[string]interface{}{
						"requests": map[string]interface{}{
							"memory": fmt.Sprintf("%dMi", memoryMB),
							"cpu":    fmt.Sprintf("%d", cpu),
						},
					},
					"devices": map[string]interface{}{
						"disks": []interface{}{
							map[string]interface{}{
								"name": "containerdisk",
								"disk": map[string]interface{}{
									"bus": "virtio",
								},
							},
						},
					},
				},
				"volumes": []interface{}{
					map[string]interface{}{
						"name": "containerdisk",
						"containerDisk": map[string]interface{}{
							"image": image,
						},
					},
				},
			},
		},
	}

	_, err := c.dyn.Resource(vmiGVR).Namespace(namespace).Create(ctx, vmi, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create VMI: %w", err)
	}
	return nil
}

type VMIStatus struct {
	Phase string
	IP    string
}

func (c *Client) GetVMIStatus(ctx context.Context, name, namespace string) (*VMIStatus, error) {
	obj, err := c.dyn.Resource(vmiGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get VMI: %w", err)
	}

	phase, _, _ := unstructured.NestedString(obj.Object, "status", "phase")

	var ip string
	interfaces, found, _ := unstructured.NestedSlice(obj.Object, "status", "interfaces")
	if found && len(interfaces) > 0 {
		if iface, ok := interfaces[0].(map[string]interface{}); ok {
			ip, _, _ = unstructured.NestedString(iface, "ipAddress")
		}
	}

	return &VMIStatus{Phase: phase, IP: ip}, nil
}

func (c *Client) DeleteVMI(ctx context.Context, name, namespace string) error {
	err := c.dyn.Resource(vmiGVR).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("delete VMI: %w", err)
	}
	return nil
}
