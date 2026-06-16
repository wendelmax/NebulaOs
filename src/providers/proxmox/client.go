package proxmox

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type pveResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors json.RawMessage `json:"errors,omitempty"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	ticket     string
	csrfToken  string
	username   string
	password   string
	apiToken   string
}

func NewClient(baseURL string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Transport: tr},
	}
}

func (c *Client) WithCredentials(username, password string) *Client {
	c.username = username
	c.password = password
	return c
}

func (c *Client) WithAPIToken(token string) *Client {
	c.apiToken = token
	return c
}

func (c *Client) WithTLSInsecureSkipVerify(skip bool) *Client {
	c.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skip},
	}
	return c
}

func (c *Client) Login(ctx context.Context) error {
	form := url.Values{
		"username": {c.username},
		"password": {c.password},
	}
	body := []byte(form.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/access/ticket", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("proxmox login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("proxmox login do: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("proxmox login read: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("proxmox login failed (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var pveResp pveResponse
	if err := json.Unmarshal(respBytes, &pveResp); err != nil {
		return fmt.Errorf("proxmox login decode: %w", err)
	}

	var authData struct {
		Ticket            string `json:"ticket"`
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
	}
	if err := json.Unmarshal(pveResp.Data, &authData); err != nil {
		return fmt.Errorf("proxmox login decode data: %w", err)
	}

	c.ticket = authData.Ticket
	c.csrfToken = authData.CSRFPreventionToken
	log.Printf("[proxmox] authenticated as %s", c.username)
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, body []byte, result interface{}) error {
	return c.doWithRetry(ctx, method, path, body, result, true)
}

func (c *Client) doWithRetry(ctx context.Context, method, path string, body []byte, result interface{}, retry bool) error {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("proxmox request: %w", err)
	}

	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.apiToken != "" {
		req.Header.Set("Authorization", "PVEAPIToken="+c.apiToken)
	} else if c.ticket != "" {
		req.Header.Set("Cookie", "PVEAuthCookie="+c.ticket)
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
			req.Header.Set("CSRFPreventionToken", c.csrfToken)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("proxmox do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized && retry && c.username != "" && c.password != "" {
		resp.Body.Close()
		if err := c.Login(ctx); err != nil {
			return fmt.Errorf("proxmox auth retry: %w", err)
		}
		return c.doWithRetry(ctx, method, path, body, result, false)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("proxmox read: %w", err)
	}

	if resp.StatusCode >= 400 {
		var pveErr pveResponse
		if json.Unmarshal(respBytes, &pveErr) == nil && len(pveErr.Errors) > 0 {
			return fmt.Errorf("proxmox API error (status %d): %s", resp.StatusCode, string(pveErr.Errors))
		}
		return fmt.Errorf("proxmox API error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var pveResp pveResponse
	if err := json.Unmarshal(respBytes, &pveResp); err != nil {
		return fmt.Errorf("proxmox decode: %w", err)
	}

	if result != nil && len(pveResp.Data) > 0 {
		if err := json.Unmarshal(pveResp.Data, result); err != nil {
			return fmt.Errorf("proxmox decode data: %w", err)
		}
	}

	return nil
}

func (c *Client) GetVersion(ctx context.Context) (string, error) {
	var result struct {
		Version string `json:"version"`
	}
	if err := c.do(ctx, http.MethodGet, "/version", nil, &result); err != nil {
		return "", err
	}
	return result.Version, nil
}

func (c *Client) GetNodes(ctx context.Context) ([]string, error) {
	var nodes []struct {
		Node string `json:"node"`
	}
	if err := c.do(ctx, http.MethodGet, "/nodes", nil, &nodes); err != nil {
		return nil, err
	}
	names := make([]string, len(nodes))
	for i, n := range nodes {
		names[i] = n.Node
	}
	return names, nil
}

func (c *Client) GetNextID(ctx context.Context) (int, error) {
	var vmid string
	if err := c.do(ctx, http.MethodGet, "/cluster/nextid", nil, &vmid); err != nil {
		return 0, err
	}
	var id int
	if _, err := fmt.Sscanf(vmid, "%d", &id); err != nil {
		return 0, fmt.Errorf("parse vmid %q: %w", vmid, err)
	}
	return id, nil
}

func (c *Client) CreateVM(ctx context.Context, node string, params map[string]interface{}) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal vm spec: %w", err)
	}
	path := fmt.Sprintf("/nodes/%s/qemu", url.PathEscape(node))
	return c.do(ctx, http.MethodPost, path, body, nil)
}

func (c *Client) DeleteVM(ctx context.Context, node string, vmid int) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d", url.PathEscape(node), vmid)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) GetVMStatus(ctx context.Context, node string, vmid int) (string, error) {
	var result struct {
		Status string `json:"status"`
	}
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/current", url.PathEscape(node), vmid)
	if err := c.do(ctx, http.MethodGet, path, nil, &result); err != nil {
		return "", err
	}
	return result.Status, nil
}

func (c *Client) CreateLXC(ctx context.Context, node string, params map[string]interface{}) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal lxc spec: %w", err)
	}
	path := fmt.Sprintf("/nodes/%s/lxc", url.PathEscape(node))
	return c.do(ctx, http.MethodPost, path, body, nil)
}

func (c *Client) GetLXCStatus(ctx context.Context, node string, vmid int) (string, error) {
	var result struct {
		Status string `json:"status"`
	}
	path := fmt.Sprintf("/nodes/%s/lxc/%d/status/current", url.PathEscape(node), vmid)
	if err := c.do(ctx, http.MethodGet, path, nil, &result); err != nil {
		return "", err
	}
	return result.Status, nil
}

func (c *Client) DeleteLXC(ctx context.Context, node string, vmid int) error {
	path := fmt.Sprintf("/nodes/%s/lxc/%d", url.PathEscape(node), vmid)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) GetStorageContent(ctx context.Context, node, storage string) ([]StorageContent, error) {
	var contents []StorageContent
	path := fmt.Sprintf("/nodes/%s/storage/%s/content", url.PathEscape(node), url.PathEscape(storage))
	if err := c.do(ctx, http.MethodGet, path, nil, &contents); err != nil {
		return nil, err
	}
	return contents, nil
}

func (c *Client) GetNodeStorages(ctx context.Context, node string) ([]string, error) {
	var storages []struct {
		Storage string `json:"storage"`
	}
	path := fmt.Sprintf("/nodes/%s/storage", url.PathEscape(node))
	if err := c.do(ctx, http.MethodGet, path, nil, &storages); err != nil {
		return nil, err
	}
	names := make([]string, len(storages))
	for i, s := range storages {
		names[i] = s.Storage
	}
	return names, nil
}

func (c *Client) CreateNetwork(ctx context.Context, node string, params map[string]interface{}) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal network spec: %w", err)
	}
	path := fmt.Sprintf("/nodes/%s/network", url.PathEscape(node))
	return c.do(ctx, http.MethodPost, path, body, nil)
}

type StorageContent struct {
	Volume  string `json:"volid"`
	Content string `json:"content"`
	Size    int64  `json:"size,omitempty"`
	Format  string `json:"format,omitempty"`
	Name    string `json:"name,omitempty"`
}
