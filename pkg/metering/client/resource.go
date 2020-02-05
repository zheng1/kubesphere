package meteringclient

import (
	"fmt"

	"k8s.io/klog"
)

// 1. init config with billing service => controller + CRD
// 2. webhook service for k8s service => metering-apiserver
// 3. metering service for billing service => metering-apiserver
// 4. metering service for kubesphere => apiserver -> metering-apiserver

func (c *Client) CreateResource(productName, resourceId string) (string, error) {
	p, err := c.createProductionInstance(productName, resourceId)
	if err != nil {
		klog.Error(err)
		return "", err
	}
	s, err := c.createSubscription(p, resourceId)
	if err != nil {
		klog.Error(err)
		return "", err
	}
	return fmt.Sprintf("%s/%s", p, s), err
}
