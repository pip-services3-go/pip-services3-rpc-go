package clients

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pip-services3-go/pip-services3-commons-go/v3/data"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/v3/data"
)

type JsonHttpClient struct {
	retryablehttp.Client
	headers *data.StringValueMap
	url     string
}

// NewJsonHttpClient Creates new JsonHttpClient
func NewJsonHttpClient() *JsonHttpClient {
	jhc := JsonHttpClient{}
	jhc.Client = *retryablehttp.NewClient()
	return &jhc
}

func (c *JsonHttpClient) SetUrl(url string) {
	c.url = url
}

func (c *JsonHttpClient) SetHeaders(headers *cdata.StringValueMap) {
	c.headers = headers
}

func (c *JsonHttpClient) Get(route string, callback func(err error, req *http.Request, res http.Response, data interface{})) {

	req, _ := http.NewRequest("GET", c.url+route, nil)
	req.Header.Set("Accept", "application/json")
	//req.Header.Set("User-Agent", c.UserAgent)
	for k, v := range c.headers.Value() {
		req.Header.Set(k, v)
	}

	resp, respErr := c.Client.HTTPClient.Do(req)
	defer resp.Body.Close()
	if callback != nil {
		callback(respErr, req, *resp, resp.Body)
	}
}

func (c *JsonHttpClient) Head(route string, callback func(err error, req *http.Request, res http.Response, data interface{})) {
	req, _ := http.NewRequest("HEAD", c.url+route, nil)
	req.Header.Set("Accept", "application/json")
	//req.Header.Set("User-Agent", c.UserAgent)
	for k, v := range c.headers.Value() {
		req.Header.Set(k, v)
	}

	resp, respErr := c.Client.HTTPClient.Do(req)
	defer resp.Body.Close()
	if callback != nil {
		callback(respErr, req, *resp, resp.Body)
	}
}

func (c *JsonHttpClient) Post(route string, data interface{}, callback func(err error, req *http.Request, res http.Response, data interface{})) {
	req, _ := http.NewRequest("POST", c.url+route, nil)
	req.Header.Set("Accept", "application/json")
	//req.Header.Set("User-Agent", c.UserAgent)
	for k, v := range c.headers.Value() {
		req.Header.Set(k, v)
	}

	resp, respErr := c.Client.HTTPClient.Do(req)
	defer resp.Body.Close()
	if callback != nil {
		callback(respErr, req, *resp, resp.Body)
	}
}

func (c *JsonHttpClient) Put(route string, data interface{}, callback func(err error, req *http.Request, res http.Response, data interface{})) {
	req, _ := http.NewRequest("PUT", c.url+route, nil)
	req.Header.Set("Accept", "application/json")
	//req.Header.Set("User-Agent", c.UserAgent)
	for k, v := range c.headers.Value() {
		req.Header.Set(k, v)
	}

	resp, respErr := c.Client.HTTPClient.Do(req)
	defer resp.Body.Close()
	if callback != nil {
		callback(respErr, req, *resp, resp.Body)
	}
}

func (c *JsonHttpClient) Del(route string, callback func(err error, req *http.Request, res http.Response, data interface{})) {
	req, _ := http.NewRequest("DELETE", c.url+route, nil)
	req.Header.Set("Accept", "application/json")
	//req.Header.Set("User-Agent", c.UserAgent)
	for k, v := range c.headers.Value() {
		req.Header.Set(k, v)
	}

	resp, respErr := c.Client.HTTPClient.Do(req)
	defer resp.Body.Close()
	if callback != nil {
		callback(respErr, req, *resp, resp.Body)
	}
}
