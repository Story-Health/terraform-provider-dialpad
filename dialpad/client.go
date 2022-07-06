package dialpad

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type DialpadError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DialpadErrorWrapper struct {
	Error DialpadError `json:"error"`
}

type Client struct {
	HttpClient *http.Client
	Token      string
}

func NewClient(token string) *Client {
	client := Client{
		HttpClient: &http.Client{},
		Token:      token,
	}

	return &client
}

func (c *Client) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("https://dialpad.com/api/v2%s", path), body)
}

func (c *Client) Do(req *http.Request) ([]byte, error) {
	token := c.Token

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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
		dialpadError := &DialpadErrorWrapper{}

		err = json.Unmarshal(body, &dialpadError)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("Dialpad Error: %s", dialpadError.Error.Message))
	}

	return body, nil
}
