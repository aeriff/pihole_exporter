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

package exporter

import (
	"strings"

	"github.com/aeriff/pihole_exporter/pihole"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	namespace = "pihole"
)

var (
	domainsBeingBlocked = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "domains_being_blocked"),
		"Domains being blocked.",
		nil, nil,
	)
	dnsQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dns_queries_today"),
		"DNS Queries today.",
		nil, nil,
	)
	adsBlocked = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ads_blocked_today"),
		"Ads blocked today.",
		nil, nil,
	)
	adsPercentage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ads_percentage_today"),
		"Ads percentage today.",
		nil, nil,
	)
	topQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_queries"),
		"Top queries.",
		[]string{"domain"}, nil,
	)
	topAds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_ads"),
		"Top Ads.",
		[]string{"domain"}, nil,
	)
	topSources = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_sources"),
		"Top sources.",
		[]string{"client"}, nil,
	)
	queryTypes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "query_types"),
		"DNS Query types.",
		[]string{"type"}, nil,
	)
	forwardDestinations = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "forward_destinations"),
		"DNS Query forward destinations.",
		[]string{"target"}, nil,
	)
	lastUpdated = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "gravity_last_updated"),
		"Timestamp of last Gravity update.",
		nil, nil,
	)
	enabled = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "enabled"),
		"Whether Pi-hole is enabled or not.",
		nil, nil,
	)
)

// Exporter collects Pihole stats from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	Pihole *pihole.Client
}

// NewExporter returns an initialized Exporter.
func NewExporter(endpoint string) (*Exporter, error) {
	pihole, err := pihole.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	return &Exporter{
		Pihole: pihole,
	}, nil
}

// Describe describes all the metrics ever exported by the Pihole exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- domainsBeingBlocked
	ch <- dnsQueries
	ch <- adsBlocked
	ch <- adsPercentage
	ch <- topQueries
	ch <- topAds
	ch <- topSources
	ch <- queryTypes
	ch <- forwardDestinations
	ch <- lastUpdated
	ch <- enabled
}

// Collect the stats from channel and delivers them as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	resp, err := e.Pihole.GetMetrics()
	if err != nil {
		log.Errorf("Error collecting metrics: %s", err.Error())
		return
	}

	switch resp.Status {
	case "enabled":
		storeMetric(ch, 1, enabled)
	default:
		storeMetric(ch, 0, enabled)
	}

	storeMetric(ch, resp.DomainsBeingBlocked, domainsBeingBlocked)
	storeMetric(ch, resp.DNSQueriesToday, dnsQueries)
	storeMetric(ch, resp.AdsBlockedToday, adsBlocked)
	storeMetric(ch, resp.AdsPercentageToday, adsPercentage)
	storeMetric(ch, resp.GravityLastUpdated.Absolute, lastUpdated)

	for qtype, requests := range resp.Querytypes {
		storeMetric(ch, requests, queryTypes, strings.Fields(qtype)[0])
	}
	for domain, hits := range resp.TopQueries {
		storeMetric(ch, hits, topQueries, domain)
	}
	for domain, hits := range resp.TopAds {
		storeMetric(ch, hits, topAds, domain)
	}
	for client, requests := range resp.TopSources {
		storeMetric(ch, requests, topSources, client)
	}
	for target, requests := range resp.ForwardDestinations {
		d := strings.Split(target, "|")
		storeMetric(ch, requests, forwardDestinations, d[len(d)-1])
	}
}

func storeMetric(ch chan<- prometheus.Metric, value float64, desc *prometheus.Desc, labels ...string) {
	ch <- prometheus.MustNewConstMetric(
		desc, prometheus.GaugeValue, value, labels...)
}
