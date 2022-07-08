package dialpad

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type dialpadError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type dialpadErrorWrapper struct {
	Error dialpadError `json:"error"`
}

type client struct {
	HttpClient *http.Client
	ApiKey     string
}

func NewClient(apiKey string) *client {
	client := client{
		HttpClient: &http.Client{},
		ApiKey:     apiKey,
	}

	return &client
}

func (c *client) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("https://dialpad.com/api/v2%s", path), body)
}

func (c *client) Do(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		dialpadError := &dialpadErrorWrapper{}

		err = json.Unmarshal(body, &dialpadError)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("Dialpad Error: %s", dialpadError.Error.Message))
	}

	return body, nil
}
