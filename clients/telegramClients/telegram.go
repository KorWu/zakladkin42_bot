package telegramClients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}
func (c *Client) GetUpdates(offset int, limit int) ([]Update, error) {
	params := url.Values{}
	params.Add("offset", strconv.Itoa(offset))
	params.Add("limit", strconv.Itoa(limit))

	data, err := c.DoRequest("getUpdates", params)

	if err != nil {
		return nil, err
	}

	var result UpdatesResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response: %w", err)
	}

	if !result.Ok {
		return nil, fmt.Errorf("failed to get the updates")
	}

	return result.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	params := url.Values{}
	params.Add("chat_id", strconv.Itoa(chatID))
	params.Add("text", text)

	_, err := c.DoRequest("sendMessage", params)

	if err != nil {
		return fmt.Errorf("failed to send the message: %w", err)
	}

	return nil
}

func (c *Client) DoRequest(method string, params url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("can't make request: %w", err)
	}

	req.URL.RawQuery = params.Encode()

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("can't read response: %w", err)
	}

	return body, nil
}
