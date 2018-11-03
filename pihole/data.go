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

// This the package for the Pi HOLE API
// See: https://github.com/pi-hole/AdminLTE

package pihole

// Metrics define PiHole Prometheus metrics
type Metrics struct {
	DomainsBeingBlocked float64 `json:"domains_being_blocked"`
	DNSQueriesToday     float64 `json:"dns_queries_today"`
	AdsBlockedToday     float64 `json:"ads_blocked_today"`
	AdsPercentageToday  float64 `json:"ads_percentage_today"`
	UniqueDomains       float64 `json:"unique_domains"`
	QueriesForwarded    float64 `json:"queries_forwarded"`
	QueriesCached       float64 `json:"queries_cached"`
	ClientsEverSeen     float64 `json:"clients_ever_seen"`
	UniqueClients       float64 `json:"unique_clients"`
	DNSQueriesAllTypes  float64 `json:"dns_queries_all_types"`
	ReplyNODATA         float64 `json:"reply_NODATA"`
	ReplyNXDOMAIN       float64 `json:"reply_NXDOMAIN"`
	ReplyCNAME          float64 `json:"reply_CNAME"`
	ReplyIP             float64 `json:"reply_IP"`
	Status              string  `json:"status"`
	GravityLastUpdated  struct {
		Absolute float64 `json:"absolute"`
	} `json:"gravity_last_updated"`
	TopQueries          map[string]float64 `json:"top_queries"`
	TopAds              map[string]float64 `json:"top_ads"`
	TopSources          map[string]float64 `json:"top_sources"`
	ForwardDestinations map[string]float64 `json:"forward_destinations"`
	Querytypes          map[string]float64 `json:"querytypes"`
}
