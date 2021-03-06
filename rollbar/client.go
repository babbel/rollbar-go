package rollbar

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const apiURL = "https://api.rollbar.com/api/1"

// Client represents the rollbar client.
type Client struct {
	AccessToken string
	Scheme      string
	Host        string
	BasePath    string
}

// Option adds the base url and other parameters to the client.
type Option func(*Client) error

// NewClient is a constructor.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	client := Client{
		AccessToken: apiKey,
		Scheme:      u.Scheme,
		Host:        u.Host,
		BasePath:    u.Path,
	}

	if err = client.parseOptions(opts...); err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *Client) parseOptions(opts ...Option) error {
	// Range over each options function and apply it to our API type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) get(pathComponents ...string) ([]byte, error) {
	return c.getWithQueryParams(map[string]string{}, pathComponents...)
}

func (c *Client) getWithQueryParams(queryParams map[string]string, pathComponents ...string) ([]byte, error) {
	url := c.url(true, queryParams, pathComponents...)
	return c.makeRequest("GET", url, nil)
}

func (c *Client) post(data []byte, pathComponents ...string) ([]byte, error) {
	url := c.url(false, map[string]string{}, pathComponents...)
	body := bytes.NewBuffer(data)
	return c.makeRequest("POST", url, body)
}

func (c *Client) delete(pathComponents ...string) error {
	url := c.url(true, map[string]string{}, pathComponents...)
	_, err := c.makeRequest("DELETE", url, nil)
	return err
}

func (c *Client) url(withAccessToken bool, queryMap map[string]string, pathComponents ...string) string {
	query := url.Values{}
	for key, value := range queryMap {
		query.Add(key, value)
	}

	if withAccessToken {
		query.Add("access_token", c.AccessToken)
	}

	components := append([]string{c.BasePath}, pathComponents...)
	path := strings.Join(components, "/")

	u := url.URL{
		Scheme:   c.Scheme,
		Host:     c.Host,
		Path:     path,
		RawQuery: query.Encode(),
	}

	return u.String()
}

func (c *Client) makeRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("The http status code is not '200' and is '%d'", resp.StatusCode)
	}

	return responseBody, nil
}
