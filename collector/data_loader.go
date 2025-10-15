// Copyright (c) 2025, 2026, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package collector

import (
	"fmt"
	"os"
    "errors"

	"github.com/BurntSushi/toml"
	"go.yaml.in/yaml/v2"
)

func (e *Exporter) reloadMetrics() bool {
	metricsToScrape, err := e.loadMetricsToScrape()
	if err != nil {
		e.logger.Error("failed to reload metrics; continuing with last known good metrics", "error", err)
		return false
	}

	e.metricsToScrape = metricsToScrape
	e.refreshCustomMetricsHashes()
	e.initCache()
	return true
}

func (e *Exporter) loadMetricsToScrape() (map[string]*Metric, error) {
	metricsToScrape := e.DefaultMetrics()

	if len(e.CustomMetricsFiles()) == 0 {
		e.logger.Debug("No custom metrics defined.")
		return metricsToScrape, nil
	}

	for _, _customMetrics := range e.CustomMetricsFiles() {
		metrics := &Metrics{}

		if err := loadMetricsConfig(_customMetrics, metrics); err != nil {
			return nil, fmt.Errorf("failed to load custom metrics %s: %w", _customMetrics, err)
		}

		e.logger.Info("Successfully loaded custom metrics from " + _customMetrics)
		mergeMetrics(metricsToScrape, metrics)
	}

	return metricsToScrape, nil
}

func (e *Exporter) merge(metrics *Metrics) {
	for _, metric := range metrics.Metric {
		e.metricsToScrape[metric.ID] = metric
	}
}

func mergeMetrics(dst map[string]*Metric, metrics *Metrics) {
	for _, metric := range metrics.Metric {
		dst[metric.ID] = metric
	}
}

func loadYamlMetricsConfig(_metricsFileName string, metrics *Metrics) error {
	yamlBytes, err := os.ReadFile(_metricsFileName)
	if err != nil {
		return fmt.Errorf("cannot read the metrics config %s: %w", _metricsFileName, err)
	}
	if err := yaml.Unmarshal(yamlBytes, metrics); err != nil {
		return fmt.Errorf("cannot unmarshal the metrics config %s: %w", _metricsFileName, err)
	}
	return nil
}

func loadTomlMetricsConfig(_customMetrics string, metrics *Metrics) error {
	if _, err := toml.DecodeFile(_customMetrics, metrics); err != nil {
		return fmt.Errorf("cannot read the metrics config %s: %w", _customMetrics, err)
	}
	return nil
}

func loadMetricsConfig(_customMetrics string, metrics *Metrics) error {
	if errToml := loadTomlMetricsConfig(_customMetrics, metrics); errToml != nil {
		if errYaml := loadYamlMetricsConfig(_customMetrics, metrics); errYaml != nil {
			return fmt.Errorf(
				"failed to load metrics config %s: %w",
				_customMetrics, errors.Join(errToml, errYaml),
			)
		}
	}
	metrics.normalizeIdentifiers()
	return nil
}
