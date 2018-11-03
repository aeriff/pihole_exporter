// Copyright (C) 2016-2018 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pihole

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint string) (*Client, error) {
	url, err := url.Parse(endpoint)
	if err != nil || url.Scheme != "http" {
		return nil, fmt.Errorf("Invalid Pi-hole address: %s", err)
	}

	resp, err := http.Get(fmt.Sprintf("%s/admin/api.php?version", url.String()))
	if err != nil {
		return nil, fmt.Errorf("Failed to check FTL version: %s", err)
	}

	var ftl struct {
		Version int `json:"version"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ftl); err != nil {
		return nil, err
	}

	if ftl.Version != 3 {
		return nil, fmt.Errorf("Unsupported FTL API version %d", ftl.Version)
	}

	return &Client{
		Endpoint: url.String(),
	}, nil
}

func (c *Client) GetMetrics() (*Metrics, error) {
	resp, err := http.Get(fmt.Sprintf("%s/admin/api.php?summaryRaw&topItems=20&getQueryTypes&getForwardDestinations&topClients&jsonForceObject", c.Endpoint))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metrics Metrics
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("Failed to decode FTL response: %s", err)
	}

	return &metrics, nil
}
