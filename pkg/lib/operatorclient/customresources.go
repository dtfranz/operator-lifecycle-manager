package operatorclient

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
)

// CustomResourceList represents a list of custom resource objects that will
// be returned from a List() operation.
type CustomResourceList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Items []*unstructured.Unstructured `json:"items"`
}

// GetCustomResource returns the custom resource as *unstructured.Unstructured by the given name.
func (c *Client) GetCustomResource(apiGroup, version, namespace, resourcePlural, resourceName string) (*unstructured.Unstructured, error) {
	klog.V(4).Infof("[GET CUSTOM RESOURCE]: %s:%s", namespace, resourceName)
	var object unstructured.Unstructured

	b, err := c.GetCustomResourceRaw(apiGroup, version, namespace, resourcePlural, resourceName)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &object); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CUSTOM RESOURCE: %v", err)
	}
	return &object, nil
}

// GetCustomResourceRaw returns the custom resource's raw body data by the given name.
func (c *Client) GetCustomResourceRaw(apiGroup, version, namespace, resourcePlural, resourceName string) ([]byte, error) {
	klog.V(4).Infof("[GET CUSTOM RESOURCE RAW]: %s:%s", namespace, resourceName)
	httpRestClient := c.extInterface.ApiextensionsV1beta1().RESTClient()
	uri := customResourceURI(apiGroup, version, namespace, resourcePlural, resourceName)
	klog.V(4).Infof("[GET]: %s", uri)

	return httpRestClient.Get().RequestURI(uri).DoRaw(context.TODO())
}

// CreateCustomResource creates the custom resource.
func (c *Client) CreateCustomResource(item *unstructured.Unstructured) error {
	klog.V(4).Infof("[CREATE CUSTOM RESOURCE]: %s:%s", item.GetNamespace(), item.GetName())
	kind := item.GetKind()
	namespace := item.GetNamespace()
	apiVersion := item.GetAPIVersion()
	apiGroup, version, err := parseAPIVersion(apiVersion)
	if err != nil {
		return err
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return c.CreateCustomResourceRaw(apiGroup, version, namespace, kind, data)
}

// CreateCustomResourceRaw creates the raw bytes of the custom resource.
func (c *Client) CreateCustomResourceRaw(apiGroup, version, namespace, kind string, data []byte) error {
	klog.V(4).Infof("[CREATE CUSTOM RESOURCE RAW]: %s:%s", namespace, kind)
	var statusCode int

	httpRestClient := c.extInterface.ApiextensionsV1beta1().RESTClient()
	uri := customResourceDefinitionURI(apiGroup, version, namespace, kind)
	klog.V(4).Infof("[POST]: %s", uri)
	result := httpRestClient.Post().RequestURI(uri).Body(data).Do(context.TODO())

	if result.Error() != nil {
		return result.Error()
	}

	result.StatusCode(&statusCode)
	klog.V(4).Infof("Written %s, status: %d", uri, statusCode)

	if statusCode != 201 {
		return fmt.Errorf("unexpected status code %d, expecting 201", statusCode)
	}
	return nil
}

// CreateCustomResourceRawIfNotFound creates the raw bytes of the custom resource if it doesn't exist, or Updates if it does exist.
// It also returns a boolean to indicate whether a new custom resource is created.
func (c *Client) CreateCustomResourceRawIfNotFound(apiGroup, version, namespace, kind, name string, data []byte) (bool, error) {
	klog.V(4).Infof("[CREATE CUSTOM RESOURCE RAW if not found]: %s:%s", namespace, name)
	_, err := c.GetCustomResource(apiGroup, version, namespace, kind, name)
	if err == nil {
		return false, nil
	}
	if !errors.IsNotFound(err) {
		return false, err
	}
	err = c.CreateCustomResourceRaw(apiGroup, version, namespace, kind, data)
	if apierrors.IsAlreadyExists(err) {
		if err = c.UpdateCustomResourceRaw(apiGroup, version, namespace, kind, name, data); err != nil {
			return false, err
		}
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateCustomResource updates the custom resource.
// To do an atomic update, use AtomicModifyCustomResource().
func (c *Client) UpdateCustomResource(item *unstructured.Unstructured) error {
	klog.V(4).Infof("[UPDATE CUSTOM RESOURCE]: %s:%s", item.GetNamespace(), item.GetName())
	kind := item.GetKind()
	name := item.GetName()
	namespace := item.GetNamespace()
	apiVersion := item.GetAPIVersion()
	apiGroup, version, err := parseAPIVersion(apiVersion)
	if err != nil {
		return err
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return c.UpdateCustomResourceRaw(apiGroup, version, namespace, kind, name, data)
}

// UpdateCustomResourceRaw updates the thirdparty resource with the raw data.
func (c *Client) UpdateCustomResourceRaw(apiGroup, version, namespace, resourcePlural, resourceName string, data []byte) error {
	klog.V(4).Infof("[UPDATE CUSTOM RESOURCE RAW]: %s:%s", namespace, resourceName)
	var statusCode int

	httpRestClient := c.extInterface.ApiextensionsV1beta1().RESTClient()
	uri := customResourceURI(apiGroup, version, namespace, resourcePlural, resourceName)
	klog.V(4).Infof("[PUT]: %s", uri)
	result := httpRestClient.Put().RequestURI(uri).Body(data).Do(context.TODO())

	if result.Error() != nil {
		return result.Error()
	}

	result.StatusCode(&statusCode)
	klog.V(4).Infof("Updated %s, status: %d", uri, statusCode)

	if statusCode != 200 {
		return fmt.Errorf("unexpected status code %d, expecting 200", statusCode)
	}
	return nil
}

// CreateOrUpdateCustomeResourceRaw creates the custom resource if it doesn't exist.
// If the custom resource exists, it updates the existing one.
func (c *Client) CreateOrUpdateCustomeResourceRaw(apiGroup, version, namespace, resourcePlural, resourceName string, data []byte) error {
	klog.V(4).Infof("[CREATE OR UPDATE UPDATE CUSTOM RESOURCE RAW]: %s:%s", namespace, resourceName)
	old, err := c.GetCustomResourceRaw(apiGroup, version, namespace, resourcePlural, resourceName)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		return c.CreateCustomResourceRaw(apiGroup, version, namespace, resourcePlural, data)
	}

	var oldSpec, newSpec unstructured.Unstructured
	if err := json.Unmarshal(old, &oldSpec); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &newSpec); err != nil {
		return err
	}

	// Set the resource version.
	newSpec.SetResourceVersion(oldSpec.GetResourceVersion())

	data, err = json.Marshal(&newSpec)
	if err != nil {
		return err
	}

	return c.UpdateCustomResourceRaw(apiGroup, version, namespace, resourcePlural, resourceName, data)
}

// DeleteCustomResource deletes the  with the given name.
func (c *Client) DeleteCustomResource(apiGroup, version, namespace, resourcePlural, resourceName string) error {
	klog.V(4).Infof("[DELETE CUSTOM RESOURCE]: %s:%s", namespace, resourceName)
	httpRestClient := c.extInterface.ApiextensionsV1beta1().RESTClient()
	uri := customResourceURI(apiGroup, version, namespace, resourcePlural, resourceName)

	klog.V(4).Infof("[DELETE]: %s", uri)
	_, err := httpRestClient.Delete().RequestURI(uri).DoRaw(context.TODO())
	return err
}

// CustomResourceModifier takes the custom resource object, and modifies it in-place.
type CustomResourceModifier func(*unstructured.Unstructured, interface{}) error

// AtomicModifyCustomResource gets the custom resource, modifies it and writes it back.
// If it's modified by other writers, we will retry until it succeeds.
func (c *Client) AtomicModifyCustomResource(apiGroup, version, namespace, resourcePlural, resourceName string, f CustomResourceModifier, data interface{}) error {
	klog.V(4).Infof("[ATOMIC MODIFY CUSTOM RESOURCE]: %s:%s", namespace, resourceName)
	return wait.PollInfinite(time.Second, func() (bool, error) {
		var customResource unstructured.Unstructured
		b, err := c.GetCustomResourceRaw(apiGroup, version, namespace, resourcePlural, resourceName)
		if err != nil {
			klog.Errorf("Failed to get CUSTOM RESOURCE %q, kind:%q: %v", resourceName, resourcePlural, err)
			return false, err
		}

		if err := json.Unmarshal(b, &customResource); err != nil {
			klog.Errorf("Failed to unmarshal CUSTOM RESOURCE %q, kind:%q: %v", resourceName, resourcePlural, err)
			return false, err
		}

		if err := f(&customResource, data); err != nil {
			klog.Errorf("Failed to modify the CUSTOM RESOURCE %q, kind:%q: %v", resourceName, resourcePlural, err)
			return false, err
		}

		if err := c.UpdateCustomResource(&customResource); err != nil {
			if errors.IsConflict(err) {
				klog.Errorf("Failed to update CUSTOM RESOURCE %q, kind:%q: %v, will retry", resourceName, resourcePlural, err)
				return false, nil
			}
			klog.Errorf("Failed to update CUSTOM RESOURCE %q, kind:%q: %v", resourceName, resourcePlural, err)
			return false, err
		}

		return true, nil
	})
}

// customResourceURI returns the URI for the thirdparty resource.
//
// Example of apiGroup: "tco.coreos.com"
// Example of version: "v1"
// Example of namespace: "default"
// Example of resourcePlural: "ChannelOperatorConfigs"
// Example of resourceName: "test-config"
func customResourceURI(apiGroup, version, namespace, resourcePlural, resourceName string) string {
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}

	return fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s/%s",
		strings.ToLower(apiGroup),
		strings.ToLower(version),
		strings.ToLower(namespace),
		strings.ToLower(resourcePlural),
		strings.ToLower(resourceName))
}

// customResourceDefinitionURI returns the URI for the CRD.
//
// Example of apiGroup: "tco.coreos.com"
// Example of version: "v1"
// Example of namespace: "default"
// Example of resourcePlural: "ChannelOperatorConfigs"
func customResourceDefinitionURI(apiGroup, version, namespace, resourcePlural string) string {
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}

	return fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s",
		strings.ToLower(apiGroup),
		strings.ToLower(version),
		strings.ToLower(namespace),
		strings.ToLower(resourcePlural))
}

// ListCustomResource lists all custom resources for the given namespace.
func (c *Client) ListCustomResource(apiGroup, version, namespace, resourcePlural string) (*CustomResourceList, error) {
	klog.V(4).Infof("LIST CUSTOM RESOURCE]: %s", resourcePlural)

	var crList CustomResourceList

	httpRestClient := c.extInterface.ApiextensionsV1beta1().RESTClient()
	uri := customResourceDefinitionURI(apiGroup, version, namespace, resourcePlural)
	klog.V(4).Infof("[GET]: %s", uri)
	bytes, err := httpRestClient.Get().RequestURI(uri).DoRaw(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to get custom resource list: %v", err)
	}

	if err := json.Unmarshal(bytes, &crList); err != nil {
		return nil, err
	}

	return &crList, nil
}

// parseAPIVersion splits "coreos.com/v1" into
// "coreos.com" and "v1".
func parseAPIVersion(apiVersion string) (apiGroup, version string, err error) {
	parts := strings.Split(apiVersion, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid format of api version %q, expecting APIGroup/Version", apiVersion)
	}
	return path.Join(parts[:len(parts)-1]...), parts[len(parts)-1], nil
}
