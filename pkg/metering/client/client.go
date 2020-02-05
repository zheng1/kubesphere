package meteringclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/klog"
	meteringconfig "kubesphere.io/kubesphere/pkg/metering/config"
)

type Client struct {
	status meteringconfig.MeteringStatus
	client *http.Client
}
type simpleJson map[string]interface{}

func GetClient() *Client {
	return &Client{
		client: &http.Client{},
		status: meteringconfig.CurrentStatus,
	}
}

func (c *Client) createProductionInstance(productName, resourceId string) (string, error) {
	product := c.status.GetProduct(productName)
	if product.ProductId == "" {
		return "", nil
	}
	params := simpleJson{
		"access_sys_id":    c.status.AccessSystemId,
		"prod_id":          product.ProductId,
		"prod_inst_id_ext": resourceId,
		"name":             resourceId,
		"description":      "",
	}
	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://139.198.121.68:9300/v1/prodinstances",
		bytes.NewReader(data))

	if err != nil {
		klog.Error(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		klog.Error(err)
		return "", err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		klog.Error(err)
		return "", err
	}

	if resp.StatusCode > http.StatusOK {
		return "", Error{resp.StatusCode, string(data)}
	}

	r := simpleJson{}
	err = json.Unmarshal(data, &r)

	if err != nil {
		klog.Error(err)
		return "", err
	}

	return r["prod_inst_id"].(string), nil
}

func (c *Client) createSubscription(productionInstanceId, resourceId string) (string, error) {
	now := time.Now()
	params := simpleJson{
		"prod_inst_id": productionInstanceId,
		"name":         resourceId,
		"description":  "",
		"billing_mode": "按需",
		"start_time":   FormatTime(now),
		"end_time":     FormatTime(now.AddDate(10, 0, 0)),
	}
	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://139.198.121.68:9300/v1/subscriptions",
		bytes.NewReader(data))

	if err != nil {
		klog.Error(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		klog.Error(err)
		return "", err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		klog.Error(err)
		return "", err
	}

	if resp.StatusCode > http.StatusOK {
		return "", Error{resp.StatusCode, string(data)}
	}

	r := simpleJson{}
	err = json.Unmarshal(data, &r)

	if err != nil {
		klog.Error(err)
		return "", err
	}

	return r["subs_id"].(string), nil
}
